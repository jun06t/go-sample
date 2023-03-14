package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"

	pb "github.com/jun06t/go-sample/opentelemetry/proto"
	"github.com/jun06t/go-sample/opentelemetry/telemetry"
)

var tracer trace.Tracer

func main() {
	backend := os.Getenv("BACKEND_ADDR")

	tp, cleanup, err := telemetry.NewTracerProvider("otel-sample")
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()
	tracer = otel.Tracer()

	h := newHandler(backend)

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(h.alive))
	mux.Handle("/hello", telemetry.NewHTTPMiddleware(h.hello))
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
		Transport: otelhttp.NewTransport(http.DefaultTransport),
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
	_, err := h.cli.SayHello(r.Context(), req)
	if err != nil {
		log.Fatal(err)
	}

	Sleep(r.Context())

	hreq, err := http.NewRequestWithContext(r.Context(), "GET", "http://httpbin.org/delay/2", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := h.hcli.Do(hreq)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func Sleep(ctx context.Context) {
	_, span := tracer.Start(ctx, "sleep")
	defer span.End()
	time.Sleep(1 * time.Second)
}
