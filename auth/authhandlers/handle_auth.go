package authhandlers

import (
	"errors"
	"net/http"

	"github.com/YankovskiyVS/eventProject/auth/main"
)

func (s *HandleAPIServer) handleAuth(w http.ResponseWriter, r *http.Request) error {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return errors.New("missing authorization header")
	}

	claims, err := main.ValidateJWT(tokenString)
	if err != nil {
		return errors.New("invalid token")
	}

	return main.WriteJSON(w, http.StatusOK, map[string]string{
		"username": claims.Username,
		"role":     claims.Role,
	})
}
