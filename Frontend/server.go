package main

import (
	"context"
	"log"

	"github.com/faagerholm/coffee/frontend/proto"
)

type server struct {
	proto.UnimplementedFrontendServiceServer

	identity proto.IdentityServiceClient
	checkout proto.CheckoutServiceClient
}

func (s server) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	res, err := s.identity.Login(ctx, req)
	if err != nil {
		log.Printf("unable to login from gateway: %v", err)
	}
	return res, err
}

func (s server) Checkout(ctx context.Context, req *proto.PlaceOrderRequest) (*proto.PlaceOrderResponse, error) {
	return s.checkout.PlaceOrder(ctx, req)
}
