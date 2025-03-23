package eventstore

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type Store struct {
	db *sql.DB
}

func NewStore() (*Store, error) {
	//Getting all required info from docker compose file
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PGHOST_2"), os.Getenv("PGPORT_2"), os.Getenv("PGUSER_2"),
		os.Getenv("PGPASSWORD_2"), os.Getenv("PGDATABASE_2"))
	db, err := sql.Open("postgres", connStr)
	if err != nil || db.Ping() != nil {
		log.Fatalf("failed to connect to database %v", err)
		return nil, err
	}
	return &Store{
		db: db,
	}, nil
}
