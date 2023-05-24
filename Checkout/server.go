package main

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/faagerholm/coffee/checkout/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	proto.UnimplementedCheckoutServiceServer

	payment  proto.PaymentServiceClient
	email    proto.EmailServiceClient
	shipping proto.ShippingServiceClient
}

func (s server) PlaceOrder(ctx context.Context, req *proto.PlaceOrderRequest) (*proto.PlaceOrderResponse, error) {
	user := req.GetUser()
	products := req.GetProducts()

	if user.Id == nil || user.Id.Prefix != proto.Prefix_PREFIX_USER {
		return nil, status.Error(codes.InvalidArgument, "invalid user argument")
	}

	// 	// Payment
	// 	if _, err := s.payment.CreatePayment(
	// 		ctx,
	// 		&proto.CreatePaymentRequest{UserId: user.Id, Products: products},
	// 	); err != nil {
	// 		return nil, err
	// 	}

	shipping, err := s.shipping.ShipOrder(ctx, &proto.ShipOrderRequest{
		Address: &proto.Address{
			Street:  user.Address.Street,
			City:    user.Address.City,
			ZipCode: user.Address.ZipCode,
			Country: user.Address.Country,
		},
		Products: products,
	})
	if err != nil {
		return nil, err
	}

	// Send Email
	body, err := generateCheckoutBody(user, products, shipping)
	if err != nil {
		return nil, err
	}
	if _, err := s.email.SendEmail(
		ctx,
		&proto.SendEmailRequest{
			Email:   user.Email,
			Subject: "Coffee shop Oy: Order placed",
			Body:    body,
		},
	); err != nil {
		return nil, err
	}

	return &proto.PlaceOrderResponse{}, nil
}

func generateCheckoutBody(user *proto.User, products []*proto.Product, shipping *proto.ShipOrderResponse) (string, error) {
	tmpl, err := template.New("demo").Parse(`
Thank you {{.FirstName}} {{.LastName}} for your order!

Order content:
{{- range $p := .Products}}
	{{ $p.Name }} {{printf "%.2f" $p.Price }}€
{{- end}}
Shipping:
	{{.Shipping.Name}} {{printf "%.2f" .Shipping.Price}}€

Total: 	{{printf "%.2f"  .Total}}€


Your tracking ID is: {{.Shipping.TrackingId}}

Best regards,
Coffee shop Oy
`)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var total float32
	for _, product := range products {
		total += product.Price
	}
	total += shipping.Price

	data := struct {
		FirstName string
		LastName  string
		Products  []*proto.Product
		Total     float32
		Shipping  struct {
			Name       string
			Price      float32
			TrackingId string
		}
	}{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Products:  products,
		Total:     total,
		Shipping: struct {
			Name       string
			Price      float32
			TrackingId string
		}{
			Name:       "Posti shipping",
			Price:      shipping.Price,
			TrackingId: shipping.TrackingId,
		},
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
