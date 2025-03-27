package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type orderStatus string

const (
	StatusCreated  orderStatus = "created"
	StatusDone     orderStatus = "done"
	StatusCanceled orderStatus = "canceled"
)

type Order struct {
	id          uuid.UUID
	userID      int
	createdAt   time.Time
	updatedAt   time.Time
	tickets     []*Ticket
	orderStatus orderStatus
	totalPrice  float32
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
		id:          uuid.New(),
		userID:      userID,
		createdAt:   time.Now().UTC(),
		updatedAt:   time.Now().UTC(),
		tickets:     tickets,
		orderStatus: StatusCreated,
		totalPrice:  totalPrice,
	}, nil
}

func (o *Order) UserID() int {
	return o.userID
}

func (o *Order) AddTicket(t *Ticket) error {
	o.tickets = append(o.tickets, t)
	return nil
}

func calculateTotalPrice(tickets []*Ticket) float32 {
	var total float32
	for _, t := range tickets {
		total += t.price
	}
	return total
}

func (o *Order) Cancel() error {
	if o.orderStatus == StatusDone {
		return errors.New("cannot cancel completed order")
	}
	o.orderStatus = StatusCanceled
	o.updatedAt = time.Now().UTC()
	return nil
}

func (o *Order) Complete() error {
	if o.orderStatus != StatusCreated {
		return errors.New("invalid status transition")
	}
	o.orderStatus = StatusDone
	o.updatedAt = time.Now().UTC()
	return nil
}
