package postgres

import (
	"database/sql"
	"fmt"
	"time"
)

func NewPostgresConnection(host, port, user, password, dbname string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func InitializeOrderDB(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			uuid VARCHAR(36) UNIQUE NOT NULL,
			order_status VARCHAR(20) NOT NULL,
			user_id INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create orders table: %w", err)
	}
	return nil
}

func InitializeTicketDB(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS tickets (
			id SERIAL PRIMARY KEY,
			uuid VARCHAR(36) UNIQUE NOT NULL,
			event_id INTEGER NOT NULL,
			order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
			price DECIMAL(10,2) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create tickets table: %w", err)
	}
	return nil
}
