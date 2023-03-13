package main

import (
	"log"
	"net/http"
)

const (
	addr = "http://localhost:8000/hello"
)

func main() {
	cli := http.Client{
		Transport: http.DefaultTransport,
	}
	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := cli.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
}
