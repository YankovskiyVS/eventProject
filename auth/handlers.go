package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, APIError{err.Error()})
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Comntent-type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type APIServer struct {
	listenAddr string
}

type APIError struct {
	Error string
}

func NewAPISServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandleFunc(s.HandleAccount))

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) HandleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		return s.HandleCreateAccount(w, r)
	}
	return nil
}

func (s *APIServer) HandleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) HandleGetAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
