package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//scan .env file
	err_env := godotenv.Load()
	if err_env != nil {
		log.Fatal("Error loading .env file")
	}

	db_port := os.Getenv("DB_PORT")

	clientOptions := options.Client().ApplyURI(db_port)
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
}
