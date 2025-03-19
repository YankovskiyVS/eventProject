package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

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

type MongoClientControler interface {
	SignIn(*User) error
	SignUp(*User) error
	CheckAuth(*User) error
	GetAllUsers(*User) error
}

var (
	mongoClient *mongo.Client
)

func initMongo() (*mongo.Database, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("env load failed: %w", err)
	}

	dbURI := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		return nil, fmt.Errorf("connection failed to %s: %w", dbURI, err)
	}

	return client.Database(dbName), nil
}

func NewMongoUserDB(db *mongo.Database) *MongoUserDB {
	return &MongoUserDB{
		db:   db,
		coll: "users",
	}
}

func (s *MongoUserDB) SignUp(u *User) error {
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

func (s *MongoUserDB) SignIn(req AuthRequest) (string, string, error) {
	collection := s.db.Collection(s.coll)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
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

	token, err := GenerateJWT(user.Username, user.Role)
	if err != nil {
		return "", "", fmt.Errorf("token generation failed: %w", err)
	}

	return token, user.Role, nil
}

func (s *MongoUserDB) GetAllUsers(ctx context.Context) ([]*User, error) {
	cursor, err := s.db.Collection(s.coll).Find(ctx, map[string]any{})
	if err != nil {
		return nil, err
	}

	users := []*User{}
	err = cursor.All(ctx, &users)
	return users, err
}
