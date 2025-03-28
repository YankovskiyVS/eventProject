package postgres

import (
	"time"
)

type OrderPostgres struct {
	ID          uint
	UUID        string // Stores domain.Order.id as string
	OrderStatus string
	UserID      uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TicketPostgres struct {
	ID      uint
	UUID    string // Stores domain.Ticket.id as string
	EventID uint
	Price   float32
	OrderID uint // Foreign key to OrderPostgres.ID
}
