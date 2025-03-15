package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Define Server
type APIServer struct {
	listenAddr string
	event      EventDB
	consumer   KafkaConsumer
	producer   KafkaProducer
}

type APIError struct {
	Error string `json:"error"`
}

func EventAPIServer(addr string, event EventDB, consumer KafkaConsumer, producer KafkaProducer) *APIServer {
	return &APIServer{listenAddr: addr,
		event:    event,
		consumer: consumer,
		producer: producer}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	//Connect router to handlers; define methods for each handler
	router.HandleFunc("/events", makeHTTPHandleFunc(s.handleCreateEvent)).Methods("POST")
	router.HandleFunc("/events/id", makeHTTPHandleFunc(s.handleUpdateEvent)).Methods("PUT")
	router.HandleFunc("/events/id", makeHTTPHandleFunc(s.handleDeleteEvent)).Methods("DELETE")

	log.Printf("Server running on %s", s.listenAddr)
	log.Fatal(http.ListenAndServe(s.listenAddr, router))
}

// Handle custom handlers that return an error (not like usual http.HandlerFunc)
func makeHTTPHandleFunc(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

// Func to use as an output in handlers; show status message JSON format
func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
