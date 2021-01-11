package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	bindAddr = "0.0.0.0:8080"
)

func main() {
	h2s := &http2.Server{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %v, http: %v, protocol: %s", r.URL.Path, r.TLS == nil, r.Proto)
	})
	server := &http.Server{
		Addr:    bindAddr,
		Handler: h2c.NewHandler(handler, h2s),
	}
	fmt.Printf("Listening %s...\n", bindAddr)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
