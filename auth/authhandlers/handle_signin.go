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

func (s *HandleAPIServer) handleSignIn(w http.ResponseWriter, r *http.Request) error {
	var req main.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.New("invalid request body")
	}

	collection := client.Database("authdb").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user authmongodb.User
	err := collection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("invalid credentials")
		}
		return errors.New("database error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return errors.New("invalid credentials")
	}

	token, err := main.GenerateJWT(user.Username, user.Role)
	if err != nil {
		return errors.New("failed to generate token")
	}

	return main.WriteJSON(w, http.StatusOK, main.AuthResponse{
		Token: token,
		Role:  user.Role,
	})
}
