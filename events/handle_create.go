package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (s *APIServer) handleCreateEvent(w http.ResponseWriter, r *http.Request) error {
	role := r.Header.Get("X-User-Role")
	if role != "admin" {
		return errors.New("status forbidden")
	}

	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		return errors.New("invalid request")
	}

	err := db.QueryRow(`INSERT INTO event_table 
	(name, description, event_date, available_tickets, ticket_price),
	`, event.Name, event.Desc, event.Date, event.Available_tickets, event.Price)

	if err != nil {
		return errors.New("service error")
	}

	return WriteJSON(w, http.StatusCreated, map[string]string{"message": "event was created successfully"})
}
