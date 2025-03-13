package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	initDB()
	//scan .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("API_PORT")

	//Start and run the server
	server := EventAPIServer(port)
	server.Run()
}
