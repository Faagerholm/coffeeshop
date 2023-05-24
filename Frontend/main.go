package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/faagerholm/coffee/frontend/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("Starting the frontend service...")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8090))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// register server
	grpcServer := grpc.NewServer(
		// register interceptors
		grpc.ChainUnaryInterceptor(
			loggingInterceptor,
		),
	)

	proto.RegisterFrontendServiceServer(grpcServer, &server{
		identity: proto.NewIdentityServiceClient(dialGrpcServer("localhost:8081")),
		checkout: proto.NewCheckoutServiceClient(dialGrpcServer("localhost:8084")),
	})
	// reflect
	reflection.Register(grpcServer)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	// start server
	go func() {
		log.Println("Starting server on port :8090")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
		wg.Done()
	}()

	// start gateway
	go func() {
		runGateway()
		wg.Done()
	}()

	wg.Wait()
}

func dialGrpcServer(addrs string) *grpc.ClientConn {
	conn, err := grpc.Dial(addrs, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %s", err)
	}
	return conn
}

func runGateway() {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	err := proto.RegisterFrontendServiceHandlerFromEndpoint(
		context.Background(),
		mux,
		fmt.Sprintf("localhost:%d", 8090),
		opts,
	)
	if err != nil {
		log.Fatalf("failed to register: %s", err)
	}
	log.Println("Starting gateway on port :8080")
	http.ListenAndServe(":8080", mux)
}

func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (
	interface{},
	error,
) {
	log.Printf("method: %s, %v", info.FullMethod, req)
	return handler(ctx, req)
}
