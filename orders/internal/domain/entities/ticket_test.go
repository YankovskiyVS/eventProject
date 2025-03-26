package entities

import (
	"testing"

	"github.com/google/uuid"
)

func TestTickets(t *testing.T) {
	ticketsStatus := NewTicketsStatus(3, 4, 20.0, 600.0)
	tickets := NewTickets(ticketsStatus, 1234)
	if tickets.ID == (uuid.UUID{}) {
		t.Errorf("Expected tickets' ID to be set, but got %s", tickets.ID)
	}
	if ticketsStatus.Number != 2 &&
		ticketsStatus.AvailableNumber != 2 &&
		ticketsStatus.Price != 20 &&
		ticketsStatus.TotalPrice != 40 {
		t.Errorf("Expected: 3 got %v, 4 got %v, 20.0 got %v, 60.0 got %v",
			ticketsStatus.Number,
			ticketsStatus.AvailableNumber,
			ticketsStatus.Price,
			ticketsStatus.TotalPrice)
	}
	if tickets.EventID != 1234 {
		t.Errorf("Expected event ID to be 1234, got %v", tickets.EventID)
	}
}
