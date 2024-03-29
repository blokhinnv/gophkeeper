// Package errors contains different types of errors.
package errors

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// ErrUnauthorized is a predefined error for unauthorized access.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrNoUsernameProvided is a predefined error for when no username is provided.
	ErrNoUsernameProvided = errors.New("no username provided")
	// ErrUsernameIsTaken is a predefined error for when username is already taken.
	ErrUsernameIsTaken = errors.New("provided login is already taken")
	// ErrNoUsernameProvided is a predefined error for unknown collection name.
	ErrUnknownCollection = errors.New("unknown collection")
	// ErrBadCredentials is a predefined error for a case of bad credentials.
	ErrBadCredentials = errors.New("username or password is incorrect")
	// ErrRecordNotFound is a predefined error for a case when the record is not found.
	ErrRecordNotFound = errors.New("document was not found")
	// ErrNoDocuments is returned by SingleResult methods when the operation that created the SingleResult did not return any documents.
	ErrNoDocuments = mongo.ErrNoDocuments
	// ErrUsernameIsTakenMongo is a predefined mongo server error for when username is already taken.
	ErrUsernameIsTakenMongo = mongo.CommandError{Code: 11000}
)
