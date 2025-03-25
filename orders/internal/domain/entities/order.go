package entities

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID          uuid.UUID
	UserID      int
	CreatedAt   time.Time
	TotalPrice  float32
	Tickets     *Tickets
	OrderStatus OrderStatus
	EventStatus EventStatus
}

type OrderStatus struct {
	Status  string
	Created bool
	Paid    bool
}

type EventStatus struct {
	Status     string
	InProgress bool
}
