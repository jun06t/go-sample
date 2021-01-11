package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/http2"
)

const (
	loop = 1000
)

var (
	target = "http://localhost:8080"
	wg     = sync.WaitGroup{}

	cli = &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}
)

func main() {
	flag.Parse()
	if args := flag.Args(); len(args) > 0 {
		target = args[0]
	}
	wg.Add(loop)

	for i := 0; i < loop; i++ {
		go func(i int) {
			defer wg.Done()

			err := do(cli, context.Background(), target)
			if err != nil {
				log.Println("error:", err)
			}
		}(i)

		time.Sleep(50 * time.Millisecond)
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
