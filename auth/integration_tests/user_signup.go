//go:build integration

package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	testDB      *mongo.Database
	serverURL   = "http://localhost:8080"
)

// TestMain sets up the MongoDB container and runs the tests
func TestMain(m *testing.M) {
	// Start MongoDB container
	ctx := context.Background()
	mongoContainer, err := startMongoContainer(ctx)
	if err != nil {
		fmt.Printf("Failed to start MongoDB container: %v\n", err)
		os.Exit(1)
	}
	defer mongoContainer.Terminate(ctx)

	// Get MongoDB connection details
	host, err := mongoContainer.Host(ctx)
	if err != nil {
		fmt.Printf("Failed to get container host: %v\n", err)
		os.Exit(1)
	}

	port, err := mongoContainer.MappedPort(ctx, "27017")
	if err != nil {
		fmt.Printf("Failed to get container port: %v\n", err)
		os.Exit(1)
	}

	mongoURI := fmt.Sprintf("mongodb://%s:%s", host, port.Port())

	// Initialize MongoDB client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: %v\n", err)
		os.Exit(1)
	}
	mongoClient = client
	testDB = client.Database("testdb")

	// Start the microservice
	go main.main()
	time.Sleep(2 * time.Second) // Wait for the server to start

	// Run tests
	code := m.Run()

	// Clean up
	testDB.Drop(ctx)
	mongoClient.Disconnect(ctx)
	os.Exit(code)
}

// startMongoContainer starts a MongoDB container using Testcontainers
func startMongoContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp"},
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "root",
			"MONGO_INITDB_ROOT_PASSWORD": "example",
		},
		WaitingFor: wait.ForLog("Waiting for connections").WithStartupTimeout(30 * time.Second),
	}
	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func TestSignUp(t *testing.T) {
	// Define the request payload
	payload := map[string]string{
		"username": "testuser",
		"password": "testpassword",
		"role":     "user",
	}
	jsonPayload, _ := json.Marshal(payload)

	// Send a POST request to sign up
	resp, err := http.Post(serverURL+"/signup", "application/json", bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Validate the response
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response map[string]string
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)
	assert.Equal(t, "User created successfully", response["message"])
	assert.NotEmpty(t, response["id"])

	// Verify the user was created in MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user bson.M
	err = testDB.Collection("users").FindOne(ctx, bson.M{"username": "testuser"}).Decode(&user)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", user["username"])
}

func TestSignIn(t *testing.T) {
	// Create a test user in MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := testDB.Collection("users").InsertOne(ctx, bson.M{
		"username": "testuser",
		"password": "$2a$10$examplehashedpassword", // Replace with a real hashed password
		"role":     "user",
	})
	assert.NoError(t, err)

	// Define the request payload
	payload := map[string]string{
		"username": "testuser",
		"password": "testpassword",
	}
	jsonPayload, _ := json.Marshal(payload)

	// Send a POST request to sign in
	resp, err := http.Post(serverURL+"/signin", "application/json", bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Validate the response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)
	assert.NotEmpty(t, response["token"])
	assert.Equal(t, "user", response["role"])
}
