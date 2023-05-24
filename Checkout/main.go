package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/faagerholm/coffee/checkout/proto"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func main() {
	log.Println("Starting the checkout service...")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8084))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// register server
	grpcServer := grpc.NewServer()
	proto.RegisterCheckoutServiceServer(grpcServer, &server{
		payment:  proto.NewPaymentServiceClient(dialGrpcServer("localhost:8085")),
		email:    proto.NewEmailServiceClient(dialGrpcServer("localhost:8086")),
		shipping: proto.NewShippingServiceClient(dialGrpcServer("localhost:8087")),
	})

	// register reflection on gRPC server
	reflection.Register(grpcServer)

	log.Println("Starting server on port :8084")
	// listen on server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func dialGrpcServer(addrs string) *grpc.ClientConn {
	conn, err := grpc.Dial(addrs, insecure.NewCredentials())
	if err != nil {
		log.Fatalf("failed to dial: %s", err)
	}
	return conn
}

func authInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (
	interface{},
	error,
) {
	// Get the token from the context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "No metadata found")
	}

	token := ""
	log.Printf("Metadata: %v", md)
	if auth := md.Get("authorization"); len(auth) > 0 {
		token = auth[0]
	}
	if token == "" {
		return nil, status.Errorf(codes.Unauthenticated, "No token found")
	}

	// Validate the token
	valid, err := validateJWT(token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if !valid {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token")
	}

	// Continue execution of handler after ensuring a valid token
	return handler(ctx, req)
}

func validateJWT(tokenString string) (bool, error) {
	// ENV variable
	secretKey := []byte("SUPER_SECRET_KEY")

	if tokenString == "" {
		return false, nil
	}

	// Define the claims struct to match your token's structure
	type Claims struct {
		jwt.StandardClaims
		// Custom claims, if any
	}

	// Parse the token with the specified secret key
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return false, err
	}

	// Validate token signature
	if !token.Valid {
		return false, nil
	}

	// Validate token claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return false, nil
	}

	return true, nil
}
