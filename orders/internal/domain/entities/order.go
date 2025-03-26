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
	Tickets     Tickets
	OrderStatus OrderStatus
	EventStatus EventStatus
}

type OrderStatusValue string

const (
	StatusCreated  OrderStatusValue = "created"
	StatusDone     OrderStatusValue = "done"
	StatusCanceled OrderStatusValue = "canceled"
)

type OrderStatus struct {
	Status  OrderStatusValue // created || done || cancelled
	Created bool
	Paid    bool
}

type EventStatusValue string

const (
	StatusEnded      EventStatusValue = "ended"
	StatusInProgress EventStatusValue = "in_progress"
	StatusCancelled  EventStatusValue = "cancelled"
)

type EventStatus struct {
	Status     EventStatusValue // ended || in_progress || cancelled
	InProgress bool
}

// Start the factory for the OrderStatus object value
func NewOrderStatus(status OrderStatusValue, created bool, paid bool) OrderStatus {
	return OrderStatus{
		Status:  status,
		Created: created,
		Paid:    paid,
	}
}

// Start the factory for the EventStatus object value
func NewEventStatus(status EventStatusValue, inProgress bool) EventStatus {
	return EventStatus{
		Status:     status,
		InProgress: inProgress,
	}
}

// Start the factory for the order entity
func NewOrder(userId int, tickets ValidatedTickets,
	orderStatus OrderStatus, eventStatus EventStatus) *Order {
	return &Order{
		ID:          uuid.New(),
		UserID:      userId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Tickets:     tickets.Tickets,
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
