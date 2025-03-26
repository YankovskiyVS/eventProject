package entities

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewOrder(t *testing.T) {
	ticketsStatus := NewTicketsStatus(3, 4, 20.0, 60.0)
	tickets := NewTickets(ticketsStatus, 1234)
	validatedTickets, err := NewValidatedTickets(tickets)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	orderStatus := NewOrderStatus("created", true, false)
	eventStatus := NewEventStatus("in_progress", true)
	order := NewOrder(1, *validatedTickets, orderStatus, eventStatus)
	if order.ID == (uuid.UUID{}) {
		t.Errorf("Expected order's ID to be set, but got %s", order.ID)
	}
	if order.UserID != 1 {
		t.Errorf("Expected user ID to be 1, got %v", order.UserID)
	}
	if orderStatus.Status != "created" {
		t.Errorf("Expected status to be 'created' got %s", orderStatus.Status)
	}
	if orderStatus.Created != true {
		t.Errorf("Expected status created true, got %v", orderStatus.Created)
	}
	if orderStatus.Paid != false {
		t.Errorf("Expected status to be not paid, got %v", orderStatus.Paid)
	}
	if eventStatus.Status != "in_progress" {
		t.Errorf("Expected event to be in progress, got %v", eventStatus.Status)
	}
	if eventStatus.InProgress != true {
		t.Errorf("Expected in progress status to be true, got %v", eventStatus.InProgress)
	}

}
