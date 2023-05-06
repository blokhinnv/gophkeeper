// Package service provides authentication services for the gophkeeper server.
package service

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/blokhinnv/gophkeeper/internal/server/auth"
	srvErrors "github.com/blokhinnv/gophkeeper/internal/server/errors"
	"github.com/blokhinnv/gophkeeper/internal/server/models"
)

// AuthService is an interface that defines the methods to handle authentication-related operations.
type AuthService interface {
	// Register creates a new user with the specified username and hashed password.
	Register(username, password string) error
	// Login attempts to authenticate a user with the specified username and password,
	// and returns a JWT token if successful.
	Login(username, password string) (string, error)
}

// authService is an implementation of the AuthService interface.
type authService struct {
	collection     *mongo.Collection // The MongoDB collection used to store user data.
	signingKey     []byte            // The signing key used to generate JWT tokens.
	expireDuration time.Duration     // The duration for which JWT tokens are valid.
}

// NewAuthService creates a new instance of the authService struct with the specified parameters.
func NewAuthService(
	collection *mongo.Collection,
	signingKey string,
	expireDuration time.Duration,
) AuthService {
	return &authService{
		collection:     collection,
		signingKey:     []byte(signingKey),
		expireDuration: expireDuration,
	}
}

// Register creates a new user with the specified username and hashed password.
// Returns an error if the username is already taken or if there is an error.
func (t *authService) Register(username, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}
	_, err = t.collection.InsertOne(ctx, bson.D{
		{Key: "username", Value: username},
		{Key: "hashedPassword", Value: string(hashedPassword)},
	})
	if err != nil {
		return err
	}
	return nil
}

// Login attempts to authenticate a user with the specified username and password.
// Returns a JWT token if authentication is successful, or an error otherwise.
func (t *authService) Login(username, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var user models.User
	err := t.collection.FindOne(ctx, bson.D{{Key: "username", Value: username}}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", err
	} else if err != nil {
		return "", err
	}
	ok, err := auth.VerifyPassword(password, user.HashedPassword)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", srvErrors.ErrUnauthorized
	}
	tok, err := auth.GenerateJWTToken(username, t.signingKey, t.expireDuration)
	if err != nil {
		return "", err
	}
	return tok, err
}
