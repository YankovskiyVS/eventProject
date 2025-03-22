package models

import "time"

type Event struct {
	Id               uint      `json:"id"`
	Name             string    `json:"name"`
	Desc             string    `json:"description"`
	Date             time.Time `json:"event_date"`
	AvailableTickets uint      `json:"available_tickets"`
	Price            uint      `json:"price"`
	IsDel            uint      `json:"is_del"`
}

func NewEvent(Name, Desc string, Date time.Time, AvailableTickets, Price uint) *Event {
	return &Event{
		Name:             Name,
		Desc:             Desc,
		Date:             Date,
		AvailableTickets: AvailableTickets,
		Price:            Price,
	}
}

type EventListResponse struct {
	Data       []Event `json:"data"`
	Page       int     `json:"page"`
	ItemsCount int     `json:"items_count"`
	TotalItems int     `json:"total_items"`
	TotalPages int     `json:"total_pages"`
}
