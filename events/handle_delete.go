package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *APIServer) handleDeleteEvent(w http.ResponseWriter, r *http.Request) error {
	role := r.Header.Get("X-User-Role")
	if role != "admin" {
		return errors.New("status forbidden")
	}

	vars := mux.Vars(r)
	eventID, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.New("invalid id format")
	}

	if err := s.event.DeleteEvent(uint(eventID)); err != nil {
		return errors.New("internal server error")
	}
	return nil
}
