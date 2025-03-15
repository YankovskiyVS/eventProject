package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *APIServer) handleUpdateEvent(w http.ResponseWriter, r *http.Request) error {
	//Check the role of the user from request header
	role := r.Header.Get("X-User-Role")
	if role != "admin" {
		return errors.New("status forbidden")
	}

	//Decode request body
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		return errors.New("invalid request")
	}

	//Get ID from URL
	vars := mux.Vars(r)
	eventID, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.New("invalid id format")
	}

	//Connect handler with the method to DB
	if err := s.event.UpdateEvent(&event, uint(eventID)); err != nil {
		return errors.New("internal server error")
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"message": "event was created successfully"})
}
