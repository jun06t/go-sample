package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(10 * time.Second)
		_, err := io.WriteString(w, "Hello, world!\n")
		if err != nil {
			fmt.Println(err)
		}
	}

	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
