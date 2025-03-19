package main

import (
	"fmt"
)

const minPassLen = 5

type User struct {
	ID       string `bson:"_id,omitempty" json:"id"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
	Role     string `bson:"role" json:"role"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func NewUserFromRequest(req *CreateUserRequest) (*User, error) {
	if err := validateCreateUserRequest(req); err != nil {
		return nil, err
	}
	return &User{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
	}, nil
}

func validateCreateUserRequest(req *CreateUserRequest) error {
	if len(req.Username) < 3 {
		return fmt.Errorf("username is too short")
	}
	if len(req.Password) < minPassLen {
		return fmt.Errorf("password is too short")
	}
	return nil
}
