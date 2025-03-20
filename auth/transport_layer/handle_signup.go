package transportlayer

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/YankovskiyVS/eventProject/auth/database"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *APIServer) handleSignUp(w http.ResponseWriter, r *http.Request) error {
	var user database.User
	//Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return errors.New("invalid request body")
	}

	//Check that the role is valid
	if user.Role != "user" && user.Role != "admin" {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid role: must be 'user' or 'admin'",
		})
	}
	log.Printf("user is %v", &user)
	//Pass request to the DB method
	err := s.mongoUser.SignUp(&user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return WriteJSON(w, http.StatusConflict, map[string]string{
				"error": "username already exists",
			})
		}
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "registration failed" + err.Error(),
		})
	}

	return WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "User created successfully",
		"id":      user.ID,
	})
}
