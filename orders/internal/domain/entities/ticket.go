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

// Getter methods:
func (t *Ticket) ID() uuid.UUID {
	return t.id
}

func (t *Ticket) Price() float32 {
	return t.price
}

func (t *Ticket) EventID() int {
	return t.eventID
}

// Setter method:
func NewTicketWithDetails(id uuid.UUID, price float32, eventID int) *Ticket {
	return &Ticket{
		id:      id,
		price:   price,
		eventID: eventID,
	}
}
