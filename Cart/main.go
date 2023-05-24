package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/faagerholm/coffee/cart/proto"
	"github.com/faagerholm/coffee/cart/store/cache"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("Starting the cart service...")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8083))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			metricsInterceptor,
		),
	)
	proto.RegisterCartServiceServer(s, &server{
		store: cache.New(),
	})

	// reflect
	reflection.Register(s)
	log.Println("Starting server on port :8083")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// debugger interceptor
func debuggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("method: %s, error: %v", info.FullMethod, err)
	}
	log.Printf("method: %s, %v", info.FullMethod, resp)
	return resp, err
}

// metrics interceptor
func metricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	s := time.Now()
	resp, err := handler(ctx, req)

	d := time.Since(s)
	log.Printf("method: %s, duration: %v", info.FullMethod, d)
	return resp, err
}
