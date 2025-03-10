package main

import (
	"errors"
	"net/http"
)

func (s *APIServer) handleAuth(w http.ResponseWriter, r *http.Request) error {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return errors.New("missing authorization header")
	}

	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return errors.New("invalid token")
	}

	return WriteJSON(w, http.StatusOK, map[string]string{
		"username": claims.Username,
		"role":     claims.Role,
	})
}
