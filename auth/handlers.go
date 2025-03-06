package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

var client *mongo.Client

type APIServer struct {
	listenAddr string
}

type APIError struct {
	Error string
}

func NewAPISServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Comntent-type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/sign up", makeHTTPHandleFunc(s.HandleCreateAccount))
	router.HandleFunc("/sign in", makeHTTPHandleFunc(s.HandleGetAccount))

	http.ListenAndServe(s.listenAddr, router)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, APIError{err.Error()})
		}
	}
}

func init() {
	// Initialize MongoDB client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *APIServer) HandleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	// Validate role
	if user.Role != "user" && user.Role != "admin" {
		http.Error(w, "Invalid role. Must be 'user' or 'admin'", http.StatusBadRequest)
		return err
	}

	collection := client.Database("authdb").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if user already exists
	var existingUser User
	err = collection.FindOne(ctx, bson.M{"id": user.Id}).Decode(&existingUser)
	if err == nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return err
	}

	// Insert new user
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User created: %s", user.Username)
	return err
}

func (s *APIServer) HandleGetAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
