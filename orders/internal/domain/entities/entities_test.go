package entities

import (
	"testing"
)

func TestEntites(t *testing.T) {
	// Test valid data
	tickets, err := NewTicket(20.0, 1)
	if err != nil {
		t.Errorf("Expected no errors in ticket %v, but got an error: %s", tickets, err)
	}
	order, err := NewOrder(1, []*Ticket{tickets})
	if err != nil {
		t.Errorf("Expected no errors in order %v, but got an error: %s", order, err)
	}

	// Test invalid data
	invalidTickets1, err := NewTicket(0, 1)
	if err == nil {
		t.Errorf("Expected an error in ticket %v", invalidTickets1)
	}
	invalidTickets2, err := NewTicket(20.0, 0)
	if err == nil {
		t.Errorf("Expected an error in ticket %v", invalidTickets2)
	}
	invalidOrder1, err := NewOrder(0, []*Ticket{tickets})
	if err == nil {
		t.Errorf("Expected an error in order %v", invalidOrder1)
	}
	invalidOrder2, err := NewOrder(1, []*Ticket{})
	if err == nil {
		t.Errorf("Expected an error in order %v", invalidOrder2)
	}

}
