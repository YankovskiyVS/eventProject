package main

import (
	"log"
	"os"

	"github.com/YankovskiyVS/eventProject/orders/internal/infrastructure/postgres"
)

func main() {
	// Initialize databases
	orderDB, err := database.NewPostgresConnection(
		os.Getenv("PGHOST_ORDER"),
		os.Getenv("PGPORT_ORDER"),
		os.Getenv("PGUSER_ORDER"),
		os.Getenv("PGPASSWORD_ORDER"),
		os.Getenv("PGDATABASE_ORDER"),
	)
	if err != nil {
		log.Fatalf("Order DB connection failed: %v", err)
	}
	defer orderDB.Close()

	ticketDB, err := database.NewPostgresConnection(
		os.Getenv("PGHOST_TICKET"),
		os.Getenv("PGPORT_TICKET"),
		os.Getenv("PGUSER_TICKET"),
		os.Getenv("PGPASSWORD_TICKET"),
		os.Getenv("PGDATABASE_TICKET"),
	)
	if err != nil {
		log.Fatalf("Ticket DB connection failed: %v", err)
	}
	defer ticketDB.Close()

	// Initialize tables
	if err := database.InitializeOrderDB(orderDB); err != nil {
		log.Fatalf("Failed to initialize order DB: %v", err)
	}

	if err := database.InitializeTicketDB(ticketDB); err != nil {
		log.Fatalf("Failed to initialize ticket DB: %v", err)
	}

	// Create repository
	repo := postgres.NewOrderRepository(orderDB, ticketDB)
}
