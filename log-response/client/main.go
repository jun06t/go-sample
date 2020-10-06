package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Foo struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	r := io.TeeReader(resp.Body, os.Stderr)

	var foo Foo
	err = json.NewDecoder(r).Decode(&foo)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(foo)
}
