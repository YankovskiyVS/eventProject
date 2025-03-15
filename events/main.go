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

	//Connect to the posgreSQL DB
	event := NewPostgresEvent()

	// Initialize Kafka producer
	brokers := []string{"localhost:9092"}
	producer, err := NewKafkaProducer(brokers)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}

	// Initialize Kafka consumer
	consumer, err := NewKafkaConsumer(brokers)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	port := os.Getenv("API_PORT")

	//Start and run the server
	server := EventAPIServer(port, event, *consumer, *producer)
	server.Run()
}
