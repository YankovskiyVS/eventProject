package handlers

import (
	"errors"
	"net/http"

	"github.com/YankovskiyVS/eventProject/auth/authhttp"
)

func (s *authhttp.APIServer) handleAuth(w http.ResponseWriter, r *http.Request) error {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return errors.New("missing authorization header")
	}

	claims, err := authjwt.validateJWT(tokenString)
	if err != nil {
		return errors.New("invalid token")
	}

	return authhttp.writeJSON(w, http.StatusOK, map[string]string{
		"username": claims.Username,
		"role":     claims.Role,
	})
}
