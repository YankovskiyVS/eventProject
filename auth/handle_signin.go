package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (s *APIServer) handleSignIn(w http.ResponseWriter, r *http.Request) error {
	var req AuthRequest
	//Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request format"})
	}

	// Input validation
	if req.Username == "" || req.Password == "" {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Username and password required"})
	}
	//Define role and token fromm DB sign in method
	token, role, err := s.mongoUser.SignIn(req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "Invalid credentials"})
		}
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Authentication failed"})
	}

	return WriteJSON(w, http.StatusOK, AuthResponse{
		Token: token,
		Role:  role,
	})
}
