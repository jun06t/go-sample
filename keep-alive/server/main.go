package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(hello))
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func hello(w http.ResponseWriter, _ *http.Request) {
	time.Sleep(50 * time.Millisecond)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Hello World")
}
