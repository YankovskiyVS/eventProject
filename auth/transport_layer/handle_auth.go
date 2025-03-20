package transportlayer

import (
	"errors"
	"net/http"
)

func (s *APIServer) handleAuth(w http.ResponseWriter, r *http.Request) error {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return errors.New("missing authorization header")
	}

	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return errors.New("invalid token")
	}

	w.Header().Set("X-User-Role", claims.Role)
	w.Header().Set("X-User-Name", claims.Username)
	w.WriteHeader(http.StatusOK)

	return WriteJSON(w, http.StatusOK, map[string]string{
		"username": claims.Username,
		"role":     claims.Role,
	})
}
