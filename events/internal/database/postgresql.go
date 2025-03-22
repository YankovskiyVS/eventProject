package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/YankovskiyVS/eventProject/events/internal/models"
	_ "github.com/lib/pq"
)

type PostgresEvent struct {
	db *sql.DB
}

// Make an interface that has all methods for the microservice
// This intrfce is declared in the APIServer (http.go file)
type EventDB interface {
	CreateEvent(*models.Event) error
	DeleteEvent(uint) error
	UpdateEvent(*models.Event, uint) error
	GetEvent(uint) (*models.Event, error)
	ListEvents(time.Time, int, int) ([]models.Event, int, error)
}

func NewPostgresEvent() *PostgresEvent {
	//Getting all required info from docker compose file
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

	//Set connection pool params
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Successfully connected to PostgreSQL database")

	return &PostgresEvent{db: db}
}

func (s *PostgresEvent) InitDB() error {
	return s.CreateEventTable()
}

func (s *PostgresEvent) CreateEventTable() error {
	//First creating of the Data table
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

func (s *PostgresEvent) CreateEvent(e *models.Event) error {
	//Creating tnew row with all info
	query := `
	INSERT INTO event_table 
	(name, description, event_date, available_tickets, ticket_price) 
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err := s.db.Query(
		query,
		e.Name,
		e.Desc,
		e.Date,
		e.AvailableTickets,
		e.Price,
	)

	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresEvent) UpdateEvent(e *models.Event, id uint) error {
	//Updating the row
	query := `UPDATE event_table 
			SET (name = $1, description = $2, event_date = $3, 
				available_tickets = $4, ticket_price = $4)
			WHERE id = $6`

	_, err := s.db.Query(
		query,
		e.Name,
		e.Desc,
		e.Date,
		e.AvailableTickets,
		e.Price,
		id,
	)

	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresEvent) DeleteEvent(id uint) error {
	//Soft-delete the row by changing is_del col
	query := `
	UPDATE event_table
	SET is_del = 1
	WHERE id = $1
	`

	_, err := s.db.Query(query, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresEvent) GetEvent(id uint) (*models.Event, error) {
	//Getting 1 event row bu ID
	query := `
	SELECT * FROM event_table
	WHERE id = $1
	`

	var event models.Event
	err := s.db.QueryRow(query, id).Scan(
		&event.Id,
		&event.Name,
		&event.Desc,
		&event.Date,
		&event.AvailableTickets,
		&event.Price,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("event not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &event, nil
}

func (s *PostgresEvent) ListEvents(dateTo *time.Time, page, itemsCount int) ([]models.Event, int, error) {
	//Get `itemsCount` of events which start date > `date`
	query := `
			SELECT * FROM event_table 
			WHERE event_date <= $1 AND is_del = 0 
			ORDER BY event_date DESC 
			LIMIT $2 OFFSET $3
			`

	//Calculate OFFSET for pagination
	offset := (page - 1) * itemsCount

	//Execute the main query
	rows, err := s.db.Query(query, dateTo, itemsCount, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.Id,
			&event.Name,
			&event.Desc,
			&event.Date,
			&event.AvailableTickets,
			&event.Price,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("row scanning error: %w", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}

	//Get total count for pagination
	countQuery := `
		SELECT COUNT(*) 
		FROM events 
		WHERE ($1::timestamp IS NULL OR date <= $1)
	`
	var total int
	err = s.db.QueryRow(countQuery, dateTo).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count query error: %w", err)
	}

	return events, total, nil
}
