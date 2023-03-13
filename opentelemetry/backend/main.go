package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/jun06t/go-sample/opentelemetry/proto"
)

const (
	port = ":8080"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Println(in.String())
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	err = s.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
