package main

import (
	"context"
	"fmt"

	"github.com/faagerholm/coffee/cart/proto"
	"github.com/faagerholm/coffee/cart/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	proto.UnimplementedCartServiceServer

	store store.Store
}

func (s *server) AddItem(ctx context.Context, req *proto.AddItemRequest) (*proto.AddItemResponse, error) {
	item := &store.Item{
		ID:       int(req.Product.Id.Value),
		Name:     req.Product.Name,
		Price:    req.Product.Price,
		Quantity: int(req.Product.Quantity),
	}

	err := s.store.AddItem(int(req.UserId.Value), *item)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("unable to add item to cart: %w", err).Error())
	}
	return &proto.AddItemResponse{Message: "OK"}, nil
}

func (s *server) GetCart(ctx context.Context, req *proto.GetCartRequest) (*proto.GetCartResponse, error) {
	userID := int(req.UserId.Value)

	cart, err := s.store.GetCart(userID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("unable to get cart: %w", err).Error())
	}

	products := make([]*proto.Product, len(cart.Items))
	for i, item := range cart.Items {
		products[i] = &proto.Product{
			Id: &proto.ID{
				Prefix: proto.Prefix_PREFIX_PAYMENT,
				Value:  int64(item.ID),
			},
			Name:  item.Name,
			Price: item.Price,
		}
	}
	return &proto.GetCartResponse{
		Products: products,
	}, nil
}

func (s *server) EmptyCart(ctx context.Context, req *proto.EmptyCartRequest) (*proto.EmptyCartResponse, error) {
	userID := int(req.UserId.Value)

	err := s.store.EmptyCart(userID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("unable to empty cart: %w", err).Error())
	}
	return &proto.EmptyCartResponse{}, nil
}
