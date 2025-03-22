package transportlayer

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (s *APIServer) handleGetEvent(w http.ResponseWriter, r *http.Request) error {
	//Check the role of the user from request header
	role := r.Header.Get("X-User-Role")
	if role != "user" {
		return errors.New("status forbidden")
	}

	vars := mux.Vars(r)
	eventID, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.New("invalid id format")
	}

	//Connect handler with the method to DB
	event := s.event.ReadEvent(uint(eventID))
	if event == nil {
		return errors.New("internal server error")
	}

	return WriteJSON(w, http.StatusOK, event)
}

func (s *APIServer) handleDateGetEvent(w http.ResponseWriter, r *http.Request) error {
	//Check the role of the user from request header
	role := r.Header.Get("X-User-Role")
	if role != "user" {
		return errors.New("status forbidden")
	}

	//Connect handler with the method to DB
	date, num := time.Now(), 20
	events, err := s.event.GetEvents(date, num)
	if err != nil {
		return errors.New("internal server error")
	}

	return WriteJSON(w, http.StatusOK, events)
}
