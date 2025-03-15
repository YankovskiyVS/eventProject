package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *APIServer) handleUpdateEvent(w http.ResponseWriter, r *http.Request) error {
	role := r.Header.Get("X-User-Role")
	if role != "admin" {
		return errors.New("status forbidden")
	}
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		return errors.New("invalid request")
	}

	vars := mux.Vars(r)
	eventID, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.New("invalid id format")
	}

	if err := s.event.UpdateEvent(&event, uint(eventID)); err != nil {
		return errors.New("internal server error")
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"message": "event was created successfully"})
}
