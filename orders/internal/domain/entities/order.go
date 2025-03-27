package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	StatusCreated  OrderStatus = "created"
	StatusDone     OrderStatus = "done"
	StatusCanceled OrderStatus = "canceled"
)

type Order struct {
	ID          uuid.UUID
	UserID      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Tickets     []*Ticket
	OrderStatus OrderStatus
	TotalPrice  float32
}

func NewOrder(userID int, tickets []*Ticket) (*Order, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}
	if len(tickets) == 0 {
		return nil, errors.New("order must contain at least one ticket")
	}

	totalPrice := calculateTotalPrice(tickets)

	return &Order{
		ID:          uuid.New(),
		UserID:      userID,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Tickets:     tickets,
		OrderStatus: StatusCreated,
		TotalPrice:  totalPrice,
	}, nil
}

func calculateTotalPrice(tickets []*Ticket) float32 {
	var total float32
	for _, t := range tickets {
		total += t.Price
	}
	return total
}

func (o *Order) Cancel() error {
	if o.OrderStatus == StatusDone {
		return errors.New("cannot cancel completed order")
	}
	o.OrderStatus = StatusCanceled
	o.UpdatedAt = time.Now().UTC()
	return nil
}

func (o *Order) Complete() error {
	if o.OrderStatus != StatusCreated {
		return errors.New("invalid status transition")
	}
	o.OrderStatus = StatusDone
	o.UpdatedAt = time.Now().UTC()
	return nil
}
