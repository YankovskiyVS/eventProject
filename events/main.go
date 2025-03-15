package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	//scan .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	event := NewPostgresEvent()

	port := os.Getenv("API_PORT")

	//Start and run the server
	server := EventAPIServer(port, event)
	server.Run()
}
