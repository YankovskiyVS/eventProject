package entities

import (
	"errors"

	"github.com/google/uuid"
)

type Ticket struct {
	ID      uuid.UUID
	Price   float32
	EventID int
}

func NewTicket(price float32, eventID int) (*Ticket, error) {
	if eventID <= 0 {
		return nil, errors.New("invalid event ID")
	}
	if price <= 0 {
		return nil, errors.New("price must be positive")
	}

	return &Ticket{
		ID:      uuid.New(),
		Price:   price,
		EventID: eventID,
	}, nil
}
