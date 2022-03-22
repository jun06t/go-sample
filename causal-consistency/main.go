package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const (
	dbName = "hoge"
	col    = "users"
)

func main() {
	ctx := context.TODO()
	cli, err := newClient(ctx, "primary:27017,secondary1:27018,secondary2:27019")
	if err != nil {
		log.Fatal(err)
	}
	/*
		ctx, closer, err := cli.withCausalConsistencySession(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer closer()
	*/

	for i := 0; i < 100; i++ {
		data := User{
			ID:   gofakeit.UUID(),
			Name: gofakeit.Name(),
			Age:  gofakeit.Number(10, 60),
		}
		list, err := cli.readYourWrite(ctx, data)
		if err != nil {
			log.Fatal(err)
		}
		if i+1 != len(list) {
			fmt.Printf("count: %d, data: %d\n", i+1, len(list))
		}
	}
	if err := cli.clean(ctx); err != nil {
		log.Fatal(err)
	}
}

type client struct {
	cli *mongo.Client
	db  *mongo.Database
}

func newClient(ctx context.Context, addr string) (*client, error) {
	mongoURI := fmt.Sprintf("mongodb://%s/?replicaSet=rs0", addr)
	opts := options.Client().ApplyURI(mongoURI).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority())).
		SetReadConcern(readconcern.Majority()).
		SetReadPreference(readpref.Secondary())

	cli, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}
	if err := cli.Connect(ctx); err != nil {
		return nil, err
	}
	ctxwt, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := cli.Ping(ctxwt, nil); err != nil {
		return nil, err
	}

	return &client{
		cli: cli,
		db:  cli.Database(dbName),
	}, nil
}

func (c *client) withCausalConsistencySession(ctx context.Context) (context.Context, func(), error) {
	sess, err := c.cli.StartSession(options.Session().SetCausalConsistency(true))
	if err != nil {
		return nil, nil, err
	}
	closer := func() {
		sess.EndSession(context.TODO())
	}
	return mongo.NewSessionContext(ctx, sess), closer, nil
}

type User struct {
	ID   string
	Name string
	Age  int
}

func (c *client) readYourWrite(ctx context.Context, data User) ([]User, error) {
	if _, err := c.db.Collection(col).InsertOne(ctx, data); err != nil {
		return nil, err
	}

	cursor, err := c.db.Collection(col).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var res []User
	if err = cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) clean(ctx context.Context) error {
	return c.db.Collection(col).Drop(ctx)
}
