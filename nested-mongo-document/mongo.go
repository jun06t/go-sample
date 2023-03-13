package main

import (
	"context"

	"github.com/chidiwilliams/flatbson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName  = "field_test"
	colName = "users"
)

type Client struct {
	cli *mongo.Client
}

func newClient(ctx context.Context) (*Client, error) {
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &Client{client}, nil
}

func (m *Client) GetUser(ctx context.Context, id string) (User, error) {
	col := m.cli.Database(dbName).Collection(colName)
	u := User{}
	err := col.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (m *Client) GetUserWithInline(ctx context.Context, id string) (UserWithInline, error) {
	col := m.cli.Database(dbName).Collection(colName)
	u := UserWithInline{}
	err := col.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	if err != nil {
		return UserWithInline{}, err
	}
	return u, nil
}

func (m *Client) UpdateUser(ctx context.Context, u User) error {
	col := m.cli.Database(dbName).Collection(colName)
	opt := options.Update().SetUpsert(true)
	_, err := col.UpdateOne(ctx, bson.M{"_id": u.ID}, bson.M{
		"$set": u,
	}, opt)
	if err != nil {
		return err
	}
	return nil
}

func (m *Client) UpdateUserRaw(ctx context.Context, u User) error {
	col := m.cli.Database(dbName).Collection(colName)
	opt := options.Update().SetUpsert(true)
	_, err := col.UpdateOne(ctx, bson.M{"_id": u.ID}, bson.M{
		"$set": bson.M{
			"age":             u.Age,
			"address.city":    u.Address.City,
			"address.zipcode": u.Address.Zipcode,
		},
	}, opt)
	if err != nil {
		return err
	}
	return nil
}

func (m *Client) UpdateUserWithInline(ctx context.Context, u UserWithInline) error {
	col := m.cli.Database(dbName).Collection(colName)
	opt := options.Update().SetUpsert(true)
	_, err := col.UpdateOne(ctx, bson.M{"_id": u.ID}, bson.M{
		"$set": u,
	}, opt)
	if err != nil {
		return err
	}
	return nil
}

func (m *Client) UpdateUserWithFlatbson(ctx context.Context, u User) error {
	col := m.cli.Database(dbName).Collection(colName)
	opt := options.Update().SetUpsert(true)
	doc, err := flatbson.Flatten(u)
	if err != nil {
		return err
	}
	_, err = col.UpdateOne(ctx, bson.M{"_id": u.ID}, bson.M{
		"$set": doc,
	}, opt)
	if err != nil {
		return err
	}
	return nil
}
