package entities

import "testing"

func TestTicketsValidation(t *testing.T) {
	// Test valid tickets
	validTicketsStatus := NewTicketsStatus(3, 4, 20.0, 60.0)
	validTickets := NewTickets(validTicketsStatus, 1234)
	if err := validTickets.validate(); err != nil {
		t.Errorf("Expected tickets to be valid, but got an error: %s", err)
	}
	// Test invalid tickets
}

func TestNewValidatedTickets(t *testing.T) {
	// Test valid tickets
	validTicketsStatus := NewTicketsStatus(3, 4, 20.0, 60.0)
	validTickets := NewTickets(validTicketsStatus, 1234)
	validatedTickets, err := NewValidatedTickets(validTickets)
	if err != nil {
		t.Errorf("Expected tickets to be valid: %v, but got an error %s", validatedTickets, err)
	}
	// Test invalid tickets
	// case 1: number of tickets > available tickets
	invalidTicketsStatus1 := NewTicketsStatus(5, 4, 20.0, 100.0)
	invalidTickets1 := NewTickets(invalidTicketsStatus1, 1234)
	validatedInvalidTickets1, err := NewValidatedTickets(invalidTickets1)
	if err == nil {
		t.Errorf("Expected tickets to be invalid: %v, but got no error", validatedInvalidTickets1)
	}
	// case 2: wrong total price
	invalidTicketsStatus2 := NewTicketsStatus(3, 4, 20.0, 20.0)
	invalidTickets2 := NewTickets(invalidTicketsStatus2, 1234)
	validatedInvalidTickets2, err := NewValidatedTickets(invalidTickets2)
	if err == nil {
		t.Errorf("Expected tickets to be invalid: %v, but got no error", validatedInvalidTickets2)
	}
	// case 3: wrong total price
	invalidTicketsStatus3 := NewTicketsStatus(3, 4, 20.0, 2.0)
	invalidTickets3 := NewTickets(invalidTicketsStatus3, 1234)
	validatedInvalidTickets3, err := NewValidatedTickets(invalidTickets3)
	if err == nil {
		t.Errorf("Expected tickets to be invalid: %v, but got no error", validatedInvalidTickets3)
	}
}
