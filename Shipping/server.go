package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"log"
	"strings"
	"time"

	"github.com/faagerholm/coffee/shipping/proto"
)

type server struct {
	proto.UnimplementedShippingServiceServer
}

func (s *server) ShipOrder(ctx context.Context, req *proto.ShipOrderRequest) (*proto.ShipOrderResponse, error) {
	log.Printf("shipping order: to: %s %s, %s, %s", req.Address.City, req.Address.Street, req.Address.ZipCode, req.Address.Country)

	return &proto.ShipOrderResponse{
		TrackingId: hash(time.Now().String()+req.Address.Street, req.Address.City, req.Address.ZipCode),
		Price:      9.90,
	}, nil
}

func hash(h ...string) string {
	// combine all strings, and hash them
	hash := md5.Sum([]byte(strings.Join(h, "-")))
	return hex.EncodeToString(hash[:])
}
