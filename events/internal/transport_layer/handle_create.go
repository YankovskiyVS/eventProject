package transportlayer

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/YankovskiyVS/eventProject/events/internal/models"
)

func (s *APIServer) handleCreateEvent(w http.ResponseWriter, r *http.Request) error {
	//Check the role of the user from request header
	role := r.Header.Get("X-User-Role")
	if role != "admin" {
		return errors.New("status forbidden")
	}

	//Decode request body
	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		return errors.New("invalid request")
	}

	//Connect handler with the method to DB
	err := s.event.CreateEvent(&event)
	if err != nil {
		return errors.New("internal server error")
	}

	return WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "Event was created successfully",
	})
}
