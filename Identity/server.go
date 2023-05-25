package main

import (
	"context"
	"time"

	"github.com/faagerholm/coffee/identity/proto"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// this should be stored in a secure location, i.e. as a
// environment variable.
var secret = []byte("SUPER_SECRET_KEY")

type server struct {
	proto.UnimplementedIdentityServiceServer
}

func (s server) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	user, err := GetByEmailAndPassword(req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	// generate JWT token
	token, err := generateJWT(user.Email + req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &proto.LoginResponse{User: user, Token: token}, nil
}

func (s server) ChangePassword(ctx context.Context, req *proto.ChangePasswordRequest) (*proto.ChangePasswordResponse, error) {
	err := changePassword(req.Email, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to change password: %v", err)
	}

	return &proto.ChangePasswordResponse{Message: "Successfully changed password!"}, nil
}

// -- JWT --
func generateJWT(in string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = in
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	claims["authorized"] = true

	return token.SignedString(secret)
}

func (s server) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	// validate token
	claims, err := validateJWT(req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	return &proto.ValidateTokenResponse{Valid: claims["authorized"].(bool)}, nil
}

func validateJWT(t string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}
