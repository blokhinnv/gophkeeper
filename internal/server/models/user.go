package models

type UserCredentials struct {
	Username string `json:"username" bson:"username" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

// User represents a user in the system. It contains a username, password, and hashed password.
type User struct {
	UserCredentials
	HashedPassword string `json:"hashedPassword" bson:"hashedPassword"`
}
