package main

import (
	"context"
	"encoding/json"

	_ "embed"

	"github.com/faagerholm/coffee/warehouse/proto"
)

type server struct {
	proto.UnimplementedWarehouseServiceServer
}

//go:embed products.json
var productJSON []byte

func (s *server) ListProducts(ctx context.Context, req *proto.Empty) (
	*proto.ListProductsResponse,
	error,
) {
	type product struct {
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float32 `json:"price"`
		Quantity    int     `json:"quantity"`
	}
	type products struct {
		Products []product `json:"products"`
	}
	var p products
	if err := json.Unmarshal(productJSON, &p); err != nil {
		return nil, err
	}

	toProtoID := func(ID int) *proto.ID {
		return &proto.ID{
			Prefix: proto.Prefix_PREFIX_PRODUCT,
			Value:  int64(ID),
		}
	}
	res := make([]*proto.Product, len(p.Products))
	for _, p := range p.Products {
		res = append(res, &proto.Product{
			Id:          toProtoID(p.ID),
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    int32(p.Quantity),
		})
	}

	return &proto.ListProductsResponse{
		Products: res,
	}, nil
}
