package main

import (
	"fmt"
	"log"
	"net"

	"github.com/faagerholm/coffee/warehouse/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("Starting warehouse service...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8082))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// register server
	grpcServer := grpc.NewServer()
	proto.RegisterWarehouseServiceServer(grpcServer, &server{})

	// register reflection on gRPC server
	reflection.Register(grpcServer)

	log.Println("starting server on port :8082")
	// listen on server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
