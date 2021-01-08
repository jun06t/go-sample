package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
)

var (
	client = http.DefaultClient
	target = "https://google.com"
)

func main() {
	flag.Parse()
	if args := flag.Args(); len(args) > 0 {
		target = args[0]
	}

	fmt.Printf("Connecting to %s...\n", target)
	err := do(context.Background(), target)
	if err != nil {
		panic(err)
	}
}

func do(ctx context.Context, url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println("Protocol:", resp.Proto)
	return nil
}
