package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cli *mongo.Client

func TestMain(m *testing.M) {
	ctx := context.TODO()
	var err error
	cli, err = mongo.Connect(
		ctx,
		options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	m.Run()
	cli.Disconnect(ctx)
}

func TestClient_UpdateUser(t *testing.T) {
	ctx := context.TODO()
	col := cli.Database(dbName).Collection(colName)

	tests := []struct {
		name     string
		fixtures func(t *testing.T)
		cleanup  func(t *testing.T)
		in       User
	}{
		{
			name: "update",
			fixtures: func(t *testing.T) {
				u := User{
					ID:   "001",
					Name: "alice",
					Age:  20,
					Address: Address{
						Country: "Japan",
						State:   "Tokyo",
						City:    "Shibuya",
						Zipcode: "150-0000",
					},
				}
				opt := options.Update().SetUpsert(true)
				_, err := col.UpdateOne(ctx, bson.M{"_id": u.ID}, bson.M{
					"$set": u,
				}, opt)
				require.NoError(t, err)
			},
			cleanup: func(t *testing.T) {
				col.Drop(ctx)
			},
			in: User{
				ID:  "001",
				Age: 25,
				Address: Address{
					City:    "Ikebukuro",
					Zipcode: "170-0000",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fixtures(t)
			defer tt.cleanup(t)

			c := &Client{cli}
			err := c.UpdateUser(ctx, tt.in)
			require.NoError(t, err)

			out, err := c.GetUser(ctx, tt.in.ID)
			require.NoError(t, err)
			b, err := json.Marshal(out)
			require.NoError(t, err)
			fmt.Println("doc: ", string(b))

			raw := make(map[string]interface{})
			err = col.FindOne(ctx, bson.M{"_id": tt.in.ID}).Decode(&raw)
			require.NoError(t, err)
			b, err = json.Marshal(raw)
			require.NoError(t, err)
			fmt.Println("raw: ", string(b))
		})
	}
}

func TestClient_UpdateUserWithInline(t *testing.T) {
	ctx := context.TODO()
	col := cli.Database(dbName).Collection(colName)

	tests := []struct {
		name     string
		fixtures func(t *testing.T)
		cleanup  func(t *testing.T)
		in       UserWithInline
	}{
		{
			name: "update",
			fixtures: func(t *testing.T) {
				u := UserWithInline{
					ID:   "001",
					Name: "alice",
					Age:  20,
					Address: Address{
						Country: "Japan",
						State:   "Tokyo",
						City:    "Shibuya",
						Zipcode: "150-0000",
					},
				}
				opt := options.Update().SetUpsert(true)
				_, err := col.UpdateOne(ctx, bson.M{"_id": u.ID}, bson.M{
					"$set": u,
				}, opt)
				require.NoError(t, err)
			},
			cleanup: func(t *testing.T) {
				col.Drop(ctx)
			},
			in: UserWithInline{
				ID:  "001",
				Age: 25,
				Address: Address{
					City:    "Ikebukuro",
					Zipcode: "170-0000",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fixtures(t)
			defer tt.cleanup(t)

			c := &Client{cli}
			err := c.UpdateUserWithInline(ctx, tt.in)
			require.NoError(t, err)

			out, err := c.GetUserWithInline(ctx, tt.in.ID)
			require.NoError(t, err)
			b, err := json.Marshal(out)
			require.NoError(t, err)
			fmt.Println("doc: ", string(b))

			raw := make(map[string]interface{})
			err = col.FindOne(ctx, bson.M{"_id": tt.in.ID}).Decode(&raw)
			require.NoError(t, err)
			b, err = json.Marshal(raw)
			require.NoError(t, err)
			fmt.Println("raw: ", string(b))
		})
	}
}

func TestClient_UpdateUserWithFlatbson(t *testing.T) {
	ctx := context.TODO()
	col := cli.Database(dbName).Collection(colName)

	tests := []struct {
		name     string
		fixtures func(t *testing.T)
		cleanup  func(t *testing.T)
		in       User
	}{
		{
			name: "update",
			fixtures: func(t *testing.T) {
				u := User{
					ID:   "001",
					Name: "alice",
					Age:  20,
					Address: Address{
						Country: "Japan",
						State:   "Tokyo",
						City:    "Shibuya",
						Zipcode: "150-0000",
					},
				}
				opt := options.Update().SetUpsert(true)
				_, err := col.UpdateOne(ctx, bson.M{"_id": u.ID}, bson.M{
					"$set": u,
				}, opt)
				require.NoError(t, err)
			},
			cleanup: func(t *testing.T) {
				col.Drop(ctx)
			},
			in: User{
				ID:  "001",
				Age: 25,
				Address: Address{
					City:    "Ikebukuro",
					Zipcode: "170-0000",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fixtures(t)
			defer tt.cleanup(t)

			c := &Client{cli}
			err := c.UpdateUserWithFlatbson(ctx, tt.in)
			require.NoError(t, err)

			out, err := c.GetUser(ctx, tt.in.ID)
			require.NoError(t, err)
			b, err := json.Marshal(out)
			require.NoError(t, err)
			fmt.Println("doc: ", string(b))

			raw := make(map[string]interface{})
			err = col.FindOne(ctx, bson.M{"_id": tt.in.ID}).Decode(&raw)
			require.NoError(t, err)
			b, err = json.Marshal(raw)
			require.NoError(t, err)
			fmt.Println("raw: ", string(b))
		})
	}
}
