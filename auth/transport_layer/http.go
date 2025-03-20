package transportlayer

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/YankovskiyVS/eventProject/auth/database"
	"github.com/gorilla/mux"
)

// Start the factory
func NewAPIServer(addr string, mongoUser *database.MongoUserDB) *APIServer {
	return &APIServer{listenAddr: addr,
		mongoUser: mongoUser}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/signup", makeHTTPHandleFunc(s.handleSignUp)).Methods("POST")
	router.HandleFunc("/signin", makeHTTPHandleFunc(s.handleSignIn)).Methods("POST")
	router.HandleFunc("/auth", makeHTTPHandleFunc(s.handleAuth)).Methods("POST")

	log.Printf("Server running on %s", s.listenAddr)
	log.Fatal(http.ListenAndServe(s.listenAddr, router))
}

func makeHTTPHandleFunc(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
