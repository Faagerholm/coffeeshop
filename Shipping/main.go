package main

import (
	"fmt"
	"log"
	"net"

	"github.com/faagerholm/coffee/shipping/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// create http listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8087))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// register server
	grpcServer := grpc.NewServer()
	proto.RegisterShippingServiceServer(grpcServer, &server{})

	// register reflection on gRPC server
	reflection.Register(grpcServer)

	// Serve server
	log.Println("starting server on port :8087")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
