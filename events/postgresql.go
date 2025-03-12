package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

type Event struct {
	Id                int       `json:"id"`
	Name              string    `json:"name"`
	Desc              string    `json:"description"`
	Date              time.Time `json:"event_date"`
	Available_tickets int       `json:"available_tickets"`
	Price             int       `json:"price"`
}

var db *sql.DB

func init() {
	var err error

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PGHOST"), os.Getenv("PGPORT"), os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"), os.Getenv("PGDATABASE"))

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS event_table (
					id SERIAL PRIMARY KEY,
					name VARCHAR(65),
					description TEXT,
					event_date DATETIME,
					available_tickets INT,
					ticket_price INT)`)
	if err != nil {
		log.Fatal(err)
	}
}
