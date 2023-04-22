// Package models provides common data structures used by the server.
package models

// Collection represents the name of a collection in the database.
type Collection string

// TextCollection, CredentialsCollection,
// BinaryCollection, and CardCollection
// are constants representing the different types
// of collections that can be used in the server.
const (
	TextCollection        = "text"
	CredentialsCollection = "credentials"
	BinaryCollection      = "binary"
	CardCollection        = "cards"
)
