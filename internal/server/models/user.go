package models

// User represents a user in the system. It contains a username, password, and hashed password.
type User struct {
	Username       string `json:"username"       bson:"username"       binding:"required"`
	Password       string `json:"password"       bson:"password"       binding:"required"`
	HashedPassword string `json:"hashedPassword" bson:"hashedPassword"`
}
