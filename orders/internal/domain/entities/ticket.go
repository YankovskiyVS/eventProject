package entities

import (
	"errors"

	"github.com/google/uuid"
)

type Tickets struct {
	ID            uuid.UUID
	TicketsStatus TicketsStatus
	EventID       int
}

type TicketsStatus struct {
	Number          int
	AvailableNumber int
	Price           float32
	TotalPrice      float32
}

// Start the factory for the TicetsStatus objest value
func NewTicketsStatus(number int, availableNumber int, price float32, totalPrice float32) TicketsStatus {
	return TicketsStatus{
		Number:          number,
		AvailableNumber: availableNumber,
		Price:           price,
		TotalPrice:      totalPrice,
	}
}

// Start the factory for the ticket entity
func NewTickets(ticketStatus TicketsStatus, eventId int) *Tickets {
	return &Tickets{
		ID:            uuid.New(),
		TicketsStatus: ticketStatus,
		EventID:       eventId,
	}
}

func (t *Tickets) validate() error {
	if t.EventID == 0 {
		return errors.New("error: cannot get the ID of the event")
	}

	if t.TicketsStatus.Number > t.TicketsStatus.AvailableNumber {
		return errors.New("error: cannot buy more tickets than available")
	}

	if t.TicketsStatus.TotalPrice < t.TicketsStatus.Price {
		return errors.New("error: total price of the tickets is less than the price of one ticket")
	}
	if float32(t.TicketsStatus.Number)*t.TicketsStatus.Price != t.TicketsStatus.TotalPrice {
		return errors.New("error: total price is not calculated right")
	}
	return nil
}

// Adds new ticket that has not been added yet
// Bug: think about TicketStatus - how to get all info from event service to the function
func (t *Tickets) AddTicket(eventId int) error {
	t.EventID = eventId
	t.TicketsStatus = TicketsStatus{}
	return t.validate()
}

func (t *Tickets) ChangeTicketsNumber(ticketsNum int) error {
	t.TicketsStatus.Number = ticketsNum
	t.TicketsStatus.TotalPrice = t.TicketsStatus.Price * float32(t.TicketsStatus.Number)
	return t.validate()
}
