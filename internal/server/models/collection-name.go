// Package models provides common data structures used by the server.
package models

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/blokhinnv/gophkeeper/internal/server/errors"
)

// CollectionName represents the name of a collection in the database.
type CollectionName string

// NewCollectionName creates a new Collection object from string and
// checks if provided value is valid.
func NewCollectionName(s string) (CollectionName, error) {
	c := CollectionName(strings.ToLower(s))
	if slices.Contains(AllowedCollectionNames, c) {
		return c, nil
	}
	return "", fmt.Errorf("%w: %v", errors.ErrUnknownCollection, s)
}

// TextCollection, CredentialsCollection,
// BinaryCollection, and CardCollection
// are constants representing the different types
// of collections that can be used in the server.
const (
	TextCollection        CollectionName = "text"
	CredentialsCollection CollectionName = "credentials"
	BinaryCollection      CollectionName = "binary"
	CardCollection        CollectionName = "cards"
)

// AllowedCollectionNames is a slice of implemented collection names.
var AllowedCollectionNames = []CollectionName{
	TextCollection,
	CredentialsCollection,
	BinaryCollection,
	CardCollection,
}
