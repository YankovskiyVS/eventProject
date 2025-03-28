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
	id          uuid.UUID
	userID      int
	createdAt   time.Time
	updatedAt   time.Time
	tickets     []*Ticket
	orderStatus OrderStatus
	totalPrice  float32
}

func NewOrder(userID int, tickets []*Ticket) (*Order, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}
	if len(tickets) == 0 {
		return nil, errors.New("order must contain at least one ticket")
	}

	totalPrice := CalculateTotalPrice(tickets)

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

// Getter methods:
func (o *Order) ID() uuid.UUID {
	return o.id
}

func (o *Order) UserID() int {
	return o.userID
}

func (o *Order) OrderStatus() OrderStatus {
	return o.orderStatus
}

func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Order) UpdatedAt() time.Time {
	return o.updatedAt
}

func (o *Order) Tickets() []*Ticket {
	return o.tickets
}

func (o *Order) AddTicket(t *Ticket) error {
	o.tickets = append(o.tickets, t)
	return nil
}

// Setter method:
func NewOrderWithDetails(
	id uuid.UUID,
	userID int,
	tickets []*Ticket,
	status string,
	createdAt time.Time,
	updatedAt time.Time,
) (*Order, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}
	if len(tickets) == 0 {
		return nil, errors.New("order must contain tickets")
	}

	return &Order{
		id:          id,
		userID:      userID,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		tickets:     tickets,
		orderStatus: OrderStatus(status),
		totalPrice:  CalculateTotalPrice(tickets),
	}, nil
}

// Busuiness logic
func CalculateTotalPrice(tickets []*Ticket) float32 {
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
