package main

import (
	"context"
	"log"
	"net"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"

	pb "github.com/jun06t/go-sample/opentelemetry/proto"
	"github.com/jun06t/go-sample/opentelemetry/telemetry"
)

const (
	port = ":8080"
)

var tracer trace.Tracer

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Println(in.String())
	Sleep(ctx)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	_, cleanup, err := telemetry.NewTracerProvider("backend")
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()
	tracer = otel.Tracer("backend")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(telemetry.NewUnaryServerInterceptor()))
	pb.RegisterGreeterServer(s, &server{})
	err = s.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}

func Sleep(ctx context.Context) {
	_, span := tracer.Start(ctx, "sleep")
	defer span.End()
	time.Sleep(200 * time.Millisecond)
}
