package main

import "github.com/faagerholm/coffee/payment/proto"

type server struct {
	proto.UnimplementedPaymentServiceServer
}
