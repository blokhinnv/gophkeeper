// Package models provides the data structures used in the application.
package models

// Metadata is a map that holds key-value pairs as additional metadata for a record.
type Metadata map[string]string

// UntypedRecord represents a record that can hold any type of data as an interface{}.
// It contains a username, data, and metadata.
type UntypedRecord struct {
	// Username represents the username of the record owner.
	Username string `json:"username" bson:"username"`
	// Data is an interface{} that can hold any type of data for the record.
	Data any `json:"data"     bson:"data"     binding:"required"`
	// Metadata is a map that can hold additional metadata for the record.
	Metadata Metadata `json:"metadata" bson:"metadata"`
}

// TextRecord represents a record that holds text data. It
// contains a username, text data, and metadata.
type TextRecord struct {
	// Username represents the username of the record owner.
	Username string
	// Data is the text data for the record.
	Data string
	// Metadata is a map that can hold additional metadata for the record.
	Metadata Metadata
}

// BinaryRecord represents a record that holds binary data.
// It contains a username, binary data, and metadata.
type BinaryRecord struct {
	// Username represents the username of the record owner.
	Username string
	// Data is the binary data for the record.
	Data []byte
	// Metadata is a map that can hold additional metadata for the record.
	Metadata Metadata
}

// Credential represents a user's login credentials.
// It contains login and password information.
type Credential struct {
	// Login represents the login information for the credential.
	Login string `validate:"required"`
	// Password represents the password information for the credential.
	Password string `validate:"required"`
}

// CredentialRecord represents a record that holds user credentials.
// It contains a username, credential data, and metadata.
type CredentialRecord struct {
	// Username represents the username of the credential owner.
	Username string
	// Data is the credential data.
	Data Credential
	// Metadata is a map that can hold additional metadata for the record.
	Metadata Metadata
}

// CardInfo represents information about a credit card.
// It contains card number, CVV, and expiration date.
type CardInfo struct {
	// CardNumber represents the card number information.
	CardNumber string `validate:"required,credit_card"`
	// CVV represents the CVV information.
	CVV string `validate:"required,len=3"`
	// ExpirationDate represents the expiration date information.
	ExpirationDate string `validate:"required,exp_date"`
}

// CardRecord represents a record that holds credit card information.
// It contains a username, card information, and metadata.
type CardRecord struct {
	// Username represents the username of the card owner.
	Username string
	// Data is the card information.
	Data CardInfo
	// Metadata is a map that can hold additional metadata for the record.
	Metadata Metadata
}
