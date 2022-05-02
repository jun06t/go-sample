package main

import (
	"context"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

var (
	endpoints = []string{"localhost:2379", "localhost:12379", "localhost:22379"}
)

const (
	lockTTL      = 10 // second
	lockResource = "/my-lock/"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints, DialTimeout: 3 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	rev, unlocker, err := Lock(context.Background(), cli, lockResource)
	if err != nil {
		log.Fatal(err)
	}
	// unlock
	defer func() {
		if err := unlocker(context.Background()); err != nil {
			log.Fatal(err)
		}
		log.Println("unlocked")
	}()
	log.Println("acquired lock rev:", rev)

	// Some function that takes a long time to complete.
	time.Sleep(5 * time.Second)
}

func Lock(ctx context.Context, cli *clientv3.Client, key string) (int64, func(context.Context) error, error) {
	ss, err := concurrency.NewSession(cli, concurrency.WithTTL(lockTTL))
	if err != nil {
		return 0, nil, err
	}
	m := concurrency.NewMutex(ss, key)
	// Orphan ends the refresh for the session lease.
	ss.Orphan()

	// acquire lock for ss
	err = m.Lock(ctx)
	// TryLock returns immediately if lock is held by another session.
	//err = m.TryLock(ctx)
	if err != nil {
		return 0, nil, err
	}

	return m.Header().Revision, func(ctx context.Context) error {
		return m.Unlock(ctx)
	}, nil
}
