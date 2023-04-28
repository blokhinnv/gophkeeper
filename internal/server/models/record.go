// Package models provides the data structures used in the application.
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Metadata is a map that holds key-value pairs as additional metadata for a record.
type Metadata map[string]string

// UntypedRecord represents a record that can hold any type of data as an interface{}.
// It contains a username, data, and metadata.
type UntypedRecord struct {
	RecordID primitive.ObjectID `json:"record_id" bson:"_id"`                         // Unique ID of a document in the DB.
	Username string             `json:"username"  bson:"username"`                    // Username represents the username of the record owner.
	Data     any                `json:"data"      bson:"data"     binding:"required"` // Data is an interface{} that can hold any type of data for the record.
	Metadata Metadata           `json:"metadata"  bson:"metadata"`                    // Metadata is a map that can hold additional metadata for the record.
}

// TextInfo is an alias for text string
type TextInfo = string

// TextRecord represents a record that holds text data. It
// contains a username, text data, and metadata.
type TextRecord struct {
	RecordID primitive.ObjectID `json:"record_id,omitempty"` // Unique ID of a document in the DB.
	Username string             `json:",omitempty"`          // Username represents the username of the record owner.
	Data     TextInfo           // Data is the text data for the record.
	Metadata Metadata           // Metadata is a map that can hold additional metadata for the record.
}

// TextInfo is an alias for text string with semantics of being base64-encoded.
type BinaryInfo = string

// BinaryRecord represents a record that holds binary data.
// It contains a username, binary data, and metadata.
type BinaryRecord struct {
	RecordID primitive.ObjectID `json:"record_id,omitempty"`                        // Unique ID of a document in the DB.
	Username string             `json:",omitempty"`                                 // Username represents the username of the record owner.
	Data     BinaryInfo         `                           validate:"credit_card"` // Data is the binary data in base64 for the record.
	Metadata Metadata           // Metadata is a map that can hold additional metadata for the record.
}

// CredentialInfo represents a user's login credentials.
// It contains login and password information.
type CredentialInfo struct {
	Login    string `validate:"required"` // Login represents the login information for the credential.
	Password string `validate:"required"` // Password represents the password information for the credential.
}

// CredentialRecord represents a record that holds user credentials.
// It contains a username, credential data, and metadata.
type CredentialRecord struct {
	RecordID primitive.ObjectID `json:"record_id,omitempty"` // Unique ID of a document in the DB.
	Username string             `json:",omitempty"`          // Username represents the username of the credential owner.
	Data     CredentialInfo     // Data is the credential data.
	Metadata Metadata           // Metadata is a map that can hold additional metadata for the record.
}

// CardInfo represents information about a credit card.
// It contains card number, CVV, and expiration date.
type CardInfo struct {
	CardNumber     string `validate:"required,credit_card"` // CardNumber represents the card number information.
	CVV            string `validate:"required,len=3"`       // CVV represents the CVV information.
	ExpirationDate string `validate:"required,exp_date"`    // ExpirationDate represents the expiration date information.
}

// CardRecord represents a record that holds credit card information.
// It contains a username, card information, and metadata.
type CardRecord struct {
	RecordID primitive.ObjectID `json:"record_id,omitempty"` // Unique ID of a document in the DB.
	Username string             `json:",omitempty"`          // Username represents the username of the card owner.
	Data     CardInfo           // Data is the card information.
	Metadata Metadata           // Metadata is a map that can hold additional metadata for the record.
}
