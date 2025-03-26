package repositories

import "github.com/YankovskiyVS/eventProject/orders/internal/domain/entities"

type TicketRepository interface {
	GetAvailableTickets(tickets *entities.ValidatedTickets)
}
