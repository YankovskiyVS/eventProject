package entities

import (
	"errors"

	"github.com/google/uuid"
)

type Ticket struct {
	id      uuid.UUID
	price   float32
	eventID int
}

func NewTicket(price float32, eventID int) (*Ticket, error) {
	if eventID <= 0 {
		return nil, errors.New("invalid event ID")
	}
	if price <= 0 {
		return nil, errors.New("price must be positive")
	}

	return &Ticket{
		id:      uuid.New(),
		price:   price,
		eventID: eventID,
	}, nil
}
