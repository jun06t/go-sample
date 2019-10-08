package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/tcnksm/go-httpstat"
)

const (
	loop = 500
)

var (
	client    = http.DefaultClient
	wg        = sync.WaitGroup{}
	customCli = &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 10 * time.Second,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          0,
			MaxIdleConnsPerHost:   loop,
			IdleConnTimeout:       60 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 20 * time.Second,
	}
)

func main() {
	url := os.Getenv("SERVER_HOST")

	wg.Add(loop)
	results := make(chan *httpstat.Result, loop)
	go reader(results)

	for i := 0; i < loop; i++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}
		result := new(httpstat.Result)
		ctx := httpstat.WithHTTPStat(req.Context(), result)
		req = req.WithContext(ctx)

		go func(i int, req *http.Request) {
			defer wg.Done()

			//err := normal(req)
			//err := readAll(req)
			err := custom(req)
			if err != nil {
				log.Println(err)
			}

			result.End(time.Now())
			fmt.Println("counter: ", i)
			results <- result
		}(i, req)

		time.Sleep(50 * time.Millisecond)
	}

	wg.Wait()
	fmt.Println("requests done")
}

func reader(results chan *httpstat.Result) {
	for result := range results {
		fmt.Printf("%+v\n", result)
	}
}

func normal(req *http.Request) error {
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func readAll(req *http.Request) error {
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	return nil
}

func custom(req *http.Request) error {
	resp, err := customCli.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	return nil
}
