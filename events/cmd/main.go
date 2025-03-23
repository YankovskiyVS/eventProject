package main

import (
	"log"
	"os"

	"github.com/YankovskiyVS/eventProject/events/internal/database"
	messagebroker "github.com/YankovskiyVS/eventProject/events/internal/message_broker"
	transportlayer "github.com/YankovskiyVS/eventProject/events/internal/transport_layer"
	"github.com/joho/godotenv"
)

func main() {
	//scan .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//Connect to the posgreSQL DB
	event := database.NewPostgresEvent()
	event.InitDB()

	// Initialize Kafka producer
	brokers := []string{"kafka:9093"}
	producer, err := messagebroker.NewKafkaProducer(brokers)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}

	// Initialize Kafka consumer
	consumer, err := messagebroker.NewKafkaConsumer(brokers)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	port := os.Getenv("API_PORT")

	//Start and run the server
	server := transportlayer.EventAPIServer(port, event, *consumer, *producer)
	server.Run()
}
