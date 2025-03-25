package entities

import "github.com/google/uuid"

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
