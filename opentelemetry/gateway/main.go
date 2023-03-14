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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"

	pb "github.com/jun06t/go-sample/opentelemetry/proto"
)

func main() {
	backend := os.Getenv("BACKEND_ADDR")

	_, cleanup, err := NewTracerProvider("otel-sample")
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	h := newHandler(backend)

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(h.alive))
	mux.Handle("/hello", otelhttp.NewHandler(http.HandlerFunc(h.hello), "hello", otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents)))
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

var tracer trace.Tracer

func Sleep(ctx context.Context) {
	_, span := tracer.Start(ctx, "sleep")
	defer span.End()
	time.Sleep(1 * time.Second)
}

func NewTracerProvider(serviceName string) (*sdktrace.TracerProvider, func(), error) {
	exporter, err := NewJaegerExporter()
	if err != nil {
		return nil, nil, err
	}

	r := NewResource(serviceName, "v1", "local")
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(r),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(1)),
	)

	otel.SetTracerProvider(tp)
	tracer = otel.Tracer("example.com/example-service")

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := tp.ForceFlush(ctx); err != nil {
			log.Print(err)
		}
		ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
		if err := tp.Shutdown(ctx2); err != nil {
			log.Print(err)
		}
		cancel()
		cancel2()
	}
	return tp, cleanup, nil
}

func NewResource(serviceName string, version string, environment string) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(version),
		attribute.String("environment", environment),
	)
}

func NewJaegerExporter() (sdktrace.SpanExporter, error) {
	// Port details: https://www.jaegertracing.io/docs/getting-started/
	endpoint := os.Getenv("EXPORTER_ENDPOINT")

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
	if err != nil {
		return nil, err
	}
	return exporter, nil
}

func NewStdoutExporter() (sdktrace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithWriter(os.Stderr),
	)
}
