// Package models provides common data structures used by the server.
package models

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/blokhinnv/gophkeeper/internal/server/errors"
)

// Collection represents the name of a collection in the database.
type Collection string

// NewCollection creates a new Collection object from string and
// checks if provided value is valid.
func NewCollection(s string) (Collection, error) {
	c := Collection(strings.ToLower(s))
	if slices.Contains(AllowedCollection, c) {
		return c, nil
	}
	return "", fmt.Errorf("%w: %v", errors.ErrUnknownCollection, s)
}

// TextCollection, CredentialsCollection,
// BinaryCollection, and CardCollection
// are constants representing the different types
// of collections that can be used in the server.
const (
	TextCollection        Collection = "text"
	CredentialsCollection Collection = "credentials"
	BinaryCollection      Collection = "binary"
	CardCollection        Collection = "cards"
)

// AllowedCollection is a slice of implemented collection names.
var AllowedCollection = []Collection{
	TextCollection,
	CredentialsCollection,
	BinaryCollection,
	CardCollection,
}
