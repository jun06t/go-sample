package telemetry

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc/filters"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
)

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

func NewUnaryServerInterceptor(opts ...otelgrpc.Option) grpc.UnaryServerInterceptor {
	opts = append(opts,
		otelgrpc.WithInterceptorFilter(
			filters.Not(filters.HealthCheck()),
		),
	)
	return otelgrpc.UnaryServerInterceptor(opts...)
}

func NewUnaryClientInterceptor(opts ...otelgrpc.Option) grpc.UnaryClientInterceptor {
	return otelgrpc.UnaryClientInterceptor(opts...)
}

func NewHTTPMiddleware(opts ...otelhttp.Option) func(http.Handler) http.Handler {
	opts = append(opts,
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
	)
	return func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, "server", opts...)
	}
}
