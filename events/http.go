package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
}

type APIError struct {
	Error string `json:"error"`
}

func EventAPIServer(addr string) *APIServer {
	return &APIServer{listenAddr: addr}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/events", makeHTTPHandleFunc(s.handleCreateEvent)).Methods("POST")
	router.HandleFunc("/events/id", makeHTTPHandleFunc(s.handleUpdateEvent)).Methods("PUT")
	router.HandleFunc("/events/id", makeHTTPHandleFunc(s.handleDeleteEvent)).Methods("DELETE")

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
