package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	//scan .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize MongoDB
	db, err := initMongo()
	if err != nil {
		log.Fatalf("MongoDB initialization failed: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Printf("MongoDB disconnect error: %v", err)
		}
	}()

	// Create MongoDB service
	mongoUserDB := NewMongoUserDB(db)

	port := os.Getenv("API_PORT")

	//Start and run the server
	server := NewAPIServer(port, mongoUserDB)
	server.Run()
}
