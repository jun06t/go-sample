package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/bigtable"
)

const (
	projectID              = "my-project"
	instanceID             = "emulator"
	tableName              = "my-table"
	columnFamilyName       = "data"
	columnQualifier        = "value"
	versionColumnQualifier = "version"
)

type Client struct {
	cli *bigtable.Client
}

func main() {
	os.Setenv("BIGTABLE_EMULATOR_HOST", "localhost:8086")
	err := initialize(columnFamilyName)
	if err != nil {
		panic(err)
	}
	cli, err := bigtable.NewClient(context.Background(), projectID, instanceID)
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	client := &Client{cli: cli}
	//err = writeOnce(client)
	//err = compareAndSwap(client)
	err = updateWhenNewVersion(client)
	if err != nil {
		panic(err)
	}
}

func writeOnce(client *Client) error {
	rowKey := "write once"
	err := client.writeOnce(rowKey, columnFamilyName, "1")
	if err != nil {
		return err
	}
	err = client.writeOnce(rowKey, columnFamilyName, "2")
	if err != nil {
		return err
	}
	err = client.read(rowKey, columnFamilyName)
	if err != nil {
		return err
	}
	return nil
}

func compareAndSwap(client *Client) error {
	rowKey := "compare and swap"
	err := client.write(rowKey, columnFamilyName, "first")
	if err != nil {
		return err
	}
	row, err := client.getRow(rowKey, columnFamilyName)
	if err != nil {
		return err
	}
	old := getValue(row, columnFamilyName, columnQualifier)

	match, err := client.compareAndSwap(rowKey, columnFamilyName, old, "second")
	if err != nil {
		return err
	}
	fmt.Println("second:", match)
	match, err = client.compareAndSwap(rowKey, columnFamilyName, old, "third")
	if err != nil {
		return err
	}
	fmt.Println("third:", match)
	err = client.read(rowKey, columnFamilyName)
	if err != nil {
		return err
	}
	return nil
}

func updateWhenNewVersion(client *Client) error {
	rowKey := "update when new version"
	err := client.updateWhenNewVersion(rowKey, columnFamilyName, "first data", "1")
	if err != nil {
		return err
	}
	err = client.updateWhenNewVersion(rowKey, columnFamilyName, "third data", "3")
	if err != nil {
		return err
	}
	err = client.updateWhenNewVersion(rowKey, columnFamilyName, "second data", "2")
	if err != nil {
		return err
	}
	err = client.read(rowKey, columnFamilyName)
	if err != nil {
		return err
	}
	return nil
}

func initialize(columnFamilyNames ...string) error {
	ctx := context.Background()
	adminClient, err := bigtable.NewAdminClient(ctx, projectID, instanceID)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	tables, err := adminClient.Tables(ctx)
	if err != nil {
		return err
	}

	if !sliceContains(tables, tableName) {
		log.Printf("Creating table %s", tableName)
		if err := adminClient.CreateTable(ctx, tableName); err != nil {
			return err
		}
	}

	tblInfo, err := adminClient.TableInfo(ctx, tableName)
	if err != nil {
		return err
	}

	for _, name := range columnFamilyNames {
		if !sliceContains(tblInfo.Families, name) {
			if err := adminClient.CreateColumnFamily(ctx, tableName, name); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Client) write(rowKey, columnFamilyName, value string) error {
	tbl := c.cli.Open(tableName)

	mut := bigtable.NewMutation()
	mut.Set(columnFamilyName, columnQualifier, bigtable.Timestamp(0), []byte(value))

	if err := tbl.Apply(context.Background(), rowKey, mut); err != nil {
		return err
	}

	return nil
}

func (c *Client) writeOnce(rowKey, columnFamilyName, value string) error {
	tbl := c.cli.Open(tableName)

	mut := bigtable.NewMutation()
	mut.Set(columnFamilyName, columnQualifier, bigtable.Timestamp(0), []byte(value))

	filter := bigtable.ChainFilters(
		bigtable.FamilyFilter(columnFamilyName),
	)
	conditionalMutation := bigtable.NewCondMutation(filter, nil, mut)

	if err := tbl.Apply(context.Background(), rowKey, conditionalMutation); err != nil {
		return err
	}

	return nil
}

func (c *Client) compareAndSwap(rowKey, columnFamilyName, old, value string) (bool, error) {
	tbl := c.cli.Open(tableName)

	mut := bigtable.NewMutation()
	mut.Set(columnFamilyName, columnQualifier, bigtable.Timestamp(0), []byte(value))

	filter := bigtable.ChainFilters(
		bigtable.FamilyFilter(columnFamilyName),
		bigtable.ColumnFilter(columnQualifier),
		bigtable.ValueFilter(old),
	)
	conditionalMutation := bigtable.NewCondMutation(filter, mut, nil)
	var match bool
	opt := bigtable.GetCondMutationResult(&match)

	if err := tbl.Apply(context.Background(), rowKey, conditionalMutation, opt); err != nil {
		return false, err
	}

	return match, nil
}

func (c *Client) updateWhenNewVersion(rowKey, columnFamilyName, value string, version string) error {
	tbl := c.cli.Open(tableName)

	mut := bigtable.NewMutation()
	mut.Set(columnFamilyName, columnQualifier, bigtable.Timestamp(0), []byte(value))
	mut.Set(columnFamilyName, versionColumnQualifier, bigtable.Timestamp(0), []byte(version))

	// if no data, set initial value
	filter := bigtable.ChainFilters(
		bigtable.FamilyFilter(columnFamilyName),
		bigtable.ColumnFilter(versionColumnQualifier),
	)
	nxCond := bigtable.NewCondMutation(filter, nil, mut)
	var match bool
	opt := bigtable.GetCondMutationResult(&match)
	if err := tbl.Apply(context.Background(), rowKey, nxCond, opt); err != nil {
		return err
	}
	if !match {
		return nil
	}

	// if data exists and version is more than current version, updates value
	filter = bigtable.ChainFilters(
		bigtable.FamilyFilter(columnFamilyName),
		bigtable.ColumnFilter(versionColumnQualifier),
		bigtable.ValueRangeFilter(nil, []byte(version)),
	)
	conditionalMutation := bigtable.NewCondMutation(filter, mut, nil)

	if err := tbl.Apply(context.Background(), rowKey, conditionalMutation); err != nil {
		return err
	}

	return nil
}

func (c *Client) getRow(rowKey, columnFamilyName string) (bigtable.Row, error) {
	tbl := c.cli.Open(tableName)
	row, err := tbl.ReadRow(context.Background(), rowKey,
		bigtable.RowFilter(
			bigtable.ChainFilters(
				bigtable.FamilyFilter(columnFamilyName),
			),
		),
	)
	return row, err
}

func (c *Client) read(rowKey, columnFamilyName string) error {
	row, err := c.getRow(rowKey, columnFamilyName)
	if err != nil {
		return err
	}

	showRow(row)
	return nil
}

func showRow(row bigtable.Row) bool {
	for _, columns := range row {
		for _, column := range columns {
			fmt.Printf("row: %s, column: %s, value: %s, timestamp: %d\n", column.Row, column.Column, string(column.Value), column.Timestamp)
		}
	}
	return true
}

func getValue(row bigtable.Row, columnFamilyName, columnQualifierName string) string {
	items := row[columnFamilyName]
	for _, item := range items {
		if item.Column == columnFamilyName+":"+columnQualifierName {
			return string(item.Value)
		}
	}
	return ""
}

// sliceContains reports whether the provided string is present in the given slice of strings.
func sliceContains(list []string, target string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}
