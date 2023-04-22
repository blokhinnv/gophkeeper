// Package auth contains functions for working with JWT tokens and generating passwords.
package auth

import "github.com/golang-jwt/jwt/v4"

// Claim is a struct containg custom fields included into JWT-token.
type Claims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
}
