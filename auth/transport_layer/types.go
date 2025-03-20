package transportlayer

import (
	"github.com/YankovskiyVS/eventProject/auth/database"
	"github.com/golang-jwt/jwt"
)

type APIServer struct {
	listenAddr string
	mongoUser  *database.MongoUserDB
}

type APIError struct {
	Error string `json:"error"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

// JWT claims structure
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}
