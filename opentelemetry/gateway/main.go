package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/grpc"

	pb "github.com/jun06t/go-sample/opentelemetry/proto"
)

func main() {
	backend := os.Getenv("BACKEND_ADDR")
	h := newHandler(backend)

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(h.alive))
	mux.Handle("/hello", http.HandlerFunc(h.hello))
	http.ListenAndServe(":8000", mux)
}

type handler struct {
	cli  pb.GreeterClient
	hcli http.Client
}

func newHandler(addr string) *handler {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	c := pb.NewGreeterClient(conn)

	hc := http.Client{
		Transport: http.DefaultTransport,
	}

	return &handler{
		cli:  c,
		hcli: hc,
	}
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
	_, err := h.cli.SayHello(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	hreq, err := http.NewRequest("GET", "http://httpbin.org/delay/2", nil)
	_, err = h.hcli.Do(hreq)
	if err != nil {
		log.Fatal(err)
	}
}
