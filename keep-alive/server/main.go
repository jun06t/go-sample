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
	svr := &http.Server{Addr: ":8080", Handler: mux}
	//svr.SetKeepAlivesEnabled(false)

	log.Fatal(svr.ListenAndServe())
}

func hello(w http.ResponseWriter, _ *http.Request) {
	time.Sleep(50 * time.Millisecond)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Hello World")
}
