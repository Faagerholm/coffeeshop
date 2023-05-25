package main

import (
	"errors"

	"github.com/faagerholm/coffee/identity/proto"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrrInvalidPassword = errors.New("invalid password")
)

// user map<password, user>
var users = map[string]*proto.User{
	"password": {
		Id: &proto.ID{
			Prefix: proto.Prefix_PREFIX_USER,
			Value:  313,
		},
		Email:     "demo@example.com",
		FirstName: "John",
		LastName:  "Doe",
	},
	"password2": {
		Id: &proto.ID{
			Prefix: proto.Prefix_PREFIX_USER,
			Value:  189,
		},
		Email:     "idis@identio.fi",
		FirstName: "Idis",
		LastName:  "Identio",
	},
}

func GetByEmailAndPassword(email, password string) (*proto.User, error) {
	user, ok := users[password]
	if !ok {
		return nil, ErrUserNotFound
	}

	if user.Email != email {
		return nil, ErrrInvalidPassword
	}

	return user, nil
}

func changePassword(email, old, new string) error {
	user, ok := users[old]
	if !ok {
		return ErrUserNotFound
	}

	if user.Email != email {
		return ErrrInvalidPassword
	}

	users[new] = user
	delete(users, old)

	return nil
}
