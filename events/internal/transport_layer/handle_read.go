package transportlayer

import (
	"errors"
	"net/http"
)

func (s *APIServer) handleGetEvent(w http.ResponseWriter, r *http.Request) error {
	//Check the role of the user from request header
	role := r.Header.Get("X-User-Role")
	if role != "user" {
		return errors.New("status forbidden")
	}

	//Connect handler with the method to DB
	err := s.event.CreateEvent(&event)
	if err != nil {
		return errors.New("internal server error")
	}

	return nil
}

func (s *APIServer) handleDateGetEvent(w http.ResponseWriter, r *http.Request) error {
	//Check the role of the user from request header
	role := r.Header.Get("X-User-Role")
	if role != "user" {
		return errors.New("status forbidden")
	}

	//Connect handler with the method to DB
	err := s.event.CreateEvent(&event)
	if err != nil {
		return errors.New("internal server error")
	}

	return nil
}
