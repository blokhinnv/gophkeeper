package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func TestValidateJWTToken(t *testing.T) {
	// Generate a sample JWT token with a known signing key
	signingKey := []byte("mySecretKey")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "blokhinnv",
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		t.Fatalf("unexpected error while generating JWT token: %v", err)
	}

	// Test case: valid token
	username, err := ValidateJWTToken(tokenString, signingKey)
	if err != nil {
		t.Fatalf("unexpected error while validating JWT token: %v", err)
	}
	if username != "blokhinnv" {
		t.Errorf("expected username to be \"blokhinnv\", but got %q", username)
	}

	// Test case: invalid token signature
	invalidSigningKey := []byte("invalidSecretKey")
	_, err = ValidateJWTToken(tokenString, invalidSigningKey)
	if err == nil {
		t.Errorf(
			"expected error while validating JWT token with invalid signing key, but got no error",
		)
	}
}

func TestGenerateJWTToken(t *testing.T) {
	// Test case: generate token with valid input
	signingKey := []byte("mySecretKey")
	username := "blokhinnv"
	expireDuration := time.Minute * 20
	tokenString, err := GenerateJWTToken(username, signingKey, expireDuration)
	if err != nil {
		t.Fatalf("unexpected error while generating JWT token: %v", err)
	}
	// Verify token signature and claims
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return signingKey, nil
		},
	)
	if err != nil {
		t.Fatalf("unexpected error while parsing JWT token: %v", err)
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		t.Fatalf("unexpected type for JWT token claims: %T", token.Claims)
	}
	if claims.Username != username {
		t.Errorf("expected username to be %q, but got %q", username, claims.Username)
	}
	if !claims.RegisteredClaims.ExpiresAt.Time.Before(time.Now().Add(expireDuration)) {
		t.Errorf(
			"expected token expiration to be within %v, but got %v",
			expireDuration,
			time.Until(claims.RegisteredClaims.ExpiresAt.Time),
		)
	}
}
