package transportlayer

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *APIServer) handleDeleteEvent(w http.ResponseWriter, r *http.Request) error {
	//Check the role of the user from request header
	role := r.Header.Get("X-User-Role")
	if role != "admin" {
		return errors.New("status forbidden")
	}

	//Get ID from URL
	vars := mux.Vars(r)
	eventID, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.New("invalid id format")
	}

	//Connect handler with the method to DB
	if err := s.event.DeleteEvent(uint(eventID)); err != nil {
		return errors.New("internal server error")
	}
	return nil
}
