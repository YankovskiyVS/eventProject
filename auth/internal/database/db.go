package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/YankovskiyVS/eventProject/auth/internal/jwt"
	"github.com/YankovskiyVS/eventProject/auth/internal/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type MongoUserDB struct {
	db   *mongo.Database
	coll string
}

func NewMongoUserDB(db *mongo.Database) *MongoUserDB {
	return &MongoUserDB{
		db:   db,
		coll: "users",
	}
}

type MongoClientControler interface {
	SignIn(*models.User) error
	SignUp(*models.User) error
	CheckAuth(*models.User) error
	GetAllUsers(*models.User) error
}

var (
	MongoClient *mongo.Client
)

var ErrInvalidCredentials = errors.New("invalid credentials")

func InitMongo() (*mongo.Database, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("env load failed: %w", err)
	}

	// Get environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Build connection string with proper scheme
	connStr := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
		url.QueryEscape(user),
		url.QueryEscape(pass),
		host,
		port,
		dbName)

	// For Docker-to-Docker communication (without auth):
	// connStr := fmt.Sprintf("mongodb://%s:%s", host, port)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	return client.Database(dbName), nil
}

func (s *MongoUserDB) SignUp(u *models.User) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("password hashing failed: %w", err)
	}
	u.Password = string(hashedPassword)

	collection := s.db.Collection(s.coll)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Increased timeout
	defer cancel()

	// Validate connection
	if err := s.db.Client().Ping(ctx, nil); err != nil {
		log.Printf("Database connection lost: %v", err)
		return fmt.Errorf("database unavailable")
	}

	// Single insert operation
	res, err := collection.InsertOne(ctx, u)
	if err != nil {
		log.Printf("database insert failed: %v", err)
		return fmt.Errorf("database insert failed: %w", err)
	}

	// Convert MongoDB ID
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		u.ID = oid.Hex()
	} else {
		return errors.New("invalid ID type")
	}

	return nil
}

func (s *MongoUserDB) SignIn(req models.AuthRequest) (string, string, error) {
	collection := s.db.Collection(s.coll)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", "", ErrInvalidCredentials
		}
		return "", "", fmt.Errorf("database error: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	); err != nil {
		return "", "", err
	}

	token, err := jwt.GenerateJWT(user.Username, user.Role)
	if err != nil {
		return "", "", fmt.Errorf("token generation failed: %w", err)
	}

	return token, user.Role, nil
}

func (s *MongoUserDB) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	cursor, err := s.db.Collection(s.coll).Find(ctx, map[string]any{})
	if err != nil {
		return nil, err
	}

	users := []*models.User{}
	err = cursor.All(ctx, &users)
	return users, err
}
