package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/tcnksm/go-httpstat"
)

const (
	url  = "http://www.google.co.jp"
	loop = 1000
)

var (
	client = http.DefaultClient
	wg     = sync.WaitGroup{}
)

func main() {
	//http.DefaultTransport.(*http.Transport).MaxIdleConns = 0
	//http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 3000

	wg.Add(loop)
	results := make(chan *httpstat.Result, 100)
	go reader(results)

	for i := 0; i < loop; i++ {
		req, _ := http.NewRequest("GET", url, nil)
		result := new(httpstat.Result)
		ctx := httpstat.WithHTTPStat(req.Context(), result)
		req = req.WithContext(ctx)

		go func(req *http.Request) {
			defer wg.Done()

			//err := normal(req)
			err := readAll(req)
			if err != nil {
				log.Println(err)
			}

			result.End(time.Now())
			results <- result
		}(req)

		time.Sleep(30 * time.Millisecond)
	}

	wg.Wait()
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
