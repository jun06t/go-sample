package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/grpc"

	pb "github.com/jun06t/go-sample/opentelemetry/proto"
)

const (
	backend = "localhost:8080"
)

func main() {
	h := newHandler(backend)

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(h.alive))
	mux.Handle("/hello", http.HandlerFunc(h.hello))
	http.ListenAndServe(":8000", mux)
}

type handler struct {
	cli grpc.Client
}

func newHandler(addr string) *handler {
	conn, err := grpc.Dial(backend, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	c := pb.NewGreeterClient(conn)

	return &handler{cli: c}
}

func (h *handler) alive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Alive")
}

func (h *handler) hello(w http.ResponseWriter, r *http.Request) {
	req := &pb.HelloRequest{
		Name: "alice",
		Age:  10,
		Man:  true,
	}
	resp, err := h.cli.SayHello(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
}
