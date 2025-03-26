package entities

import "testing"

func TestOrderValidation(t *testing.T) {
	ticketsStatus := NewTicketsStatus(3, 4, 20.0, 60.0)
	tickets := NewTickets(ticketsStatus, 1234)
	validatedTickets, err := NewValidatedTickets(tickets)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	// Test valid order
	validOrderStatus := NewOrderStatus("created", true, false)
	validEventStatus := NewEventStatus("in_progress", true)
	validOrder := NewOrder(1, *validatedTickets, validOrderStatus, validEventStatus)
	if err := validOrder.validate(); err != nil {
		t.Errorf("Expected order to be valid, but got an error: %s", err)
	}
	// Test invalid products
	invalidOrder1 := &Order{UserID: 0}
	if err := invalidOrder1.validate(); err == nil {
		t.Errorf("Expected order to be invalid but got no error")
	}
	invalidOrder2 := &Order{UserID: (-2)}
	if err := invalidOrder2.validate(); err == nil {
		t.Errorf("Expected order to be invalid but got no error")
	}
}

func TestNewValidatedOrder(t *testing.T) {
	ticketsStatus := NewTicketsStatus(3, 4, 20.0, 60.0)
	tickets := NewTickets(ticketsStatus, 1234)
	validatedTickets, err := NewValidatedTickets(tickets)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	// Test valid order
	validOrderStatus := NewOrderStatus("created", true, false)
	validEventStatus := NewEventStatus("in_progress", true)
	validOrder := NewOrder(1, *validatedTickets, validOrderStatus, validEventStatus)
	validatedOrder, err := NewValidateddOrder(validOrder)
	if err != nil {
		t.Errorf("Expected order to be valid: %v, but got error: %s", validatedOrder, err)
	}

	// Test invalid orders
	invalidOrderStatus1 := NewOrderStatus("", true, false) // Empty status
	invalidEventStatus1 := NewEventStatus("", true)        // Empty status
	invalidOrder1 := NewOrder(1, *validatedTickets, invalidOrderStatus1, invalidEventStatus1)
	validatedInvalidOrder1, err := NewValidateddOrder(invalidOrder1)
	if err == nil {
		t.Errorf("Expected order to be invalid: %v, but got no error", validatedInvalidOrder1)
	}
}
