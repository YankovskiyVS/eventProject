package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID          uuid.UUID
	UserID      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
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

// Start the factory for the OrderStatus object value
func NewOrderStatus(status string, created bool, paid bool) OrderStatus {
	return OrderStatus{
		Status:  status,
		Created: created,
		Paid:    paid,
	}
}

// Start the factory for the EventStatus object value
func NewEventStatus(status string, inProgress bool) EventStatus {
	return EventStatus{
		Status:     status,
		InProgress: inProgress,
	}
}

// Start the factory for the order entity
func NewOrder(userId int, totalPrice float32, tickets *Tickets,
	orderStatus OrderStatus, eventStatus EventStatus) *Order {
	return &Order{
		ID:          uuid.New(),
		UserID:      userId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		TotalPrice:  totalPrice,
		Tickets:     tickets,
		OrderStatus: orderStatus,
		EventStatus: eventStatus,
	}
}

// Order entity validation
func (o *Order) validate() error {
	switch {

	case o.UserID == 0:
		return errors.New("cannot get user ID")

	case o.CreatedAt.After(o.UpdatedAt):
		return errors.New("created_at must be before updated_at")
	}
	return nil
}
