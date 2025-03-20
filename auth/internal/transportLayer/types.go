package transportLayer

import "github.com/YankovskiyVS/eventProject/auth/internal/database"

type APIServer struct {
	listenAddr string
	mongoUser  *database.MongoUserDB
}

type APIError struct {
	Error string `json:"error"`
}
