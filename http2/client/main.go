package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

var (
	target = "https://google.com"

	cli1 = &http.Client{
		Transport: &http2.Transport{},
	}
	cli2 = &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 10 * time.Second,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2: false,
		},
	}
	cli3 = http.DefaultClient
)

func main() {
	flag.Parse()
	if args := flag.Args(); len(args) > 0 {
		target = args[0]
	}

	fmt.Printf("Connecting to %s...\n", target)
	err := do(cli2, context.Background(), target)
	if err != nil {
		panic(err)
	}
}

func do(cli *http.Client, ctx context.Context, url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req = req.WithContext(ctx)

	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println("Protocol:", resp.Proto)
	return nil
}
