package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/bigtable"
)

const (
	instanceID             = "jun06t-test"
	tableName              = "my-table"
	columnFamilyName       = "data"
	columnQualifier        = "value"
	versionColumnQualifier = "version"
)

type Client struct {
	writer *bigtable.Table
	reader *bigtable.Table
}

func main() {
	projectID := os.Getenv("PROJECT_ID")
	err := initialize(projectID, instanceID, tableName, columnFamilyName)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	writer, err := bigtable.NewClientWithConfig(ctx, projectID, instanceID, bigtable.ClientConfig{
		AppProfile: "writer",
	})
	if err != nil {
		panic(err)
	}
	defer writer.Close()
	reader, err := bigtable.NewClientWithConfig(ctx, projectID, instanceID, bigtable.ClientConfig{
		AppProfile: "reader",
	})
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	client := &Client{
		writer: writer.Open(tableName),
		reader: reader.Open(tableName),
	}

	var matched, unmatched int
	for i := 0; i < 100; i++ {
		match, err := writeAndRead(client)
		if err != nil {
			panic(err)
		}
		if match {
			matched++
		} else {
			unmatched++
		}
	}
	fmt.Printf("matched: %d, unmatched: %d\n", matched, unmatched)
}

func writeAndRead(client *Client) (bool, error) {
	rowKey := "key1"
	now := time.Now().UnixNano()
	in := strconv.FormatInt(now, 10)
	err := client.write(rowKey, columnFamilyName, in)
	if err != nil {
		return false, err
	}
	out, err := client.read(rowKey, columnFamilyName, false)
	if err != nil {
		return false, err
	}
	var match bool
	if in == out {
		match = true
	}
	return match, nil
}

func (c *Client) write(rowKey, columnFamilyName, value string) error {
	mut := bigtable.NewMutation()
	mut.Set(columnFamilyName, columnQualifier, bigtable.Timestamp(0), []byte(value))

	if err := c.writer.Apply(context.Background(), rowKey, mut); err != nil {
		return err
	}

	return nil
}

func (c *Client) writeOnce(rowKey, columnFamilyName, value string) error {
	mut := bigtable.NewMutation()
	mut.Set(columnFamilyName, columnQualifier, bigtable.Timestamp(0), []byte(value))

	filter := bigtable.ChainFilters(
		bigtable.FamilyFilter(columnFamilyName),
	)
	conditionalMutation := bigtable.NewCondMutation(filter, nil, mut)

	if err := c.writer.Apply(context.Background(), rowKey, conditionalMutation); err != nil {
		return err
	}

	return nil
}

func (c *Client) read(rowKey, columnFamilyName string, primary bool) (string, error) {
	tbl := c.reader
	if primary {
		tbl = c.writer
	}
	row, err := tbl.ReadRow(context.Background(), rowKey,
		bigtable.RowFilter(
			bigtable.FamilyFilter(columnFamilyName),
		),
	)
	if err != nil {
		return "", err
	}

	var out string
	for _, columns := range row {
		for _, column := range columns {
			out = string(column.Value)
		}
	}
	return out, nil
}

func showRow(row bigtable.Row) bool {
	for _, columns := range row {
		for _, column := range columns {
			fmt.Printf("row: %s, column: %s, value: %s, timestamp: %d\n", column.Row, column.Column, string(column.Value), column.Timestamp)
		}
	}
	return true
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

func initialize(projectID, instanceID, tableName string, columnFamilyNames ...string) error {
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
