package authhandlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/YankovskiyVS/eventProject/auth/authmongodb"
	"github.com/YankovskiyVS/eventProject/auth/main"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type HandleAPIServer struct {
	main.APIServer
}

var client *mongo.Client

func (s *HandleAPIServer) handleSignUp(w http.ResponseWriter, r *http.Request) error {
	var user authmongodb.User
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
	var existingUser authmongodb.User
	err = collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		return errors.New("username already exists")
	}

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return errors.New("failed to create user")
	}

	return main.WriteJSON(w, http.StatusCreated, map[string]string{"message": "User created successfully"})
}
