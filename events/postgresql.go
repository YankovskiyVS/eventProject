package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type PostgresEvent struct {
	db *sql.DB
}

type EventDB interface {
	CreateEvent(*Event) error
	DeleteEvent(uint) error
	UpdateEvent(*Event, uint) error
}

func NewPostgresEvent() *PostgresEvent {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PGHOST"), os.Getenv("PGPORT"), os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"), os.Getenv("PGDATABASE"))

	//Открыть БД
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	//Проверить соединение с БД
	if err := db.Ping(); err != nil {
		log.Fatal(err)
		return nil
	}

	return &PostgresEvent{db: db}
}

func (s *PostgresEvent) InitDB() error {
	return s.CreateEventTable()
}

func (s *PostgresEvent) CreateEventTable() error {

	queryCreate := `CREATE TABLE IF NOT EXISTS event_table (
					id SERIAL PRIMARY KEY,
					name VARCHAR(65),
					description TEXT,
					event_date DATE,
					available_tickets UINT,
					ticket_price UINT, is_del DEFAULT 0)`

	_, err := s.db.Exec(queryCreate)
	return err
}

func (s *PostgresEvent) CreateEvent(e *Event) error {
	query := `INSERT INTO event_table (name, description, event_date, available_tickets, ticket_price) 
			VALUES ($1, $2, $3, $4, $5)`

	_, err := s.db.Query(query,
		e.Name,
		e.Desc,
		e.Date,
		e.AvailableTickets,
		e.Price)

	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresEvent) UpdateEvent(e *Event, id uint) error {
	query := `UPDATE event_table 
			SET (name = $1, description = $2, event_date = $3, 
				available_tickets = $4, ticket_price = $4)
			WHERE id == $6`

	_, err := s.db.Query(query,
		e.Name,
		e.Desc,
		e.Date,
		e.AvailableTickets,
		e.Price,
		id)

	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresEvent) DeleteEvent(id uint) error {
	//Soft-delete the row by changing is_del col
	query := `UPDATE event_table SET is_del = 1 WHERE id == $1`

	_, err := s.db.Query(query, id)

	if err != nil {
		return err
	}

	return nil
}
