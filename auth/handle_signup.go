package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func (s *APIServer) handleSignUp(w http.ResponseWriter, r *http.Request) error {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return errors.New("invalid request body")
	}

	if user.Role != "user" && user.Role != "admin" {
		return errors.New("invalid role")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = string(hashedPassword)

	collection := client.Database("authdb").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check for existing user
	var existingUser User
	err = collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		return errors.New("username already exists")
	}

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return errors.New("failed to create user")
	}

	return WriteJSON(w, http.StatusCreated, map[string]string{"message": "User created successfully"})
}
