package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateJWTToken generates a JSON Web Token (JWT) for the given username and signing key,
// with an expiration time specified by the provided duration. The function returns the
// generated token as a string and an error if the token cannot be generated. The token is
// generated using the HMAC-SHA256 signing method.
//
// Parameters:
//   - username: the username to be included in the JWT claims.
//   - signingKey: the key used to sign the JWT.
//   - expireDuration: the duration after which the JWT will expire.
//
// Returns:
//   - string: the JWT token as a string.
//   - error: an error if the token cannot be generated.
func GenerateJWTToken(
	username string,
	signingKey []byte,
	expireDuration time.Duration,
) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Username: username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateJWTToken validates the provided JWT token string using the specified signing key.
// If the token is valid, the function returns the username included in the token claims as
// a string. If the token is not valid, an error is returned. The function uses the HMAC-SHA256
// signing method to validate the token.
//
// Parameters:
//   - tokenString: the JWT token string to validate.
//   - signingKey: the key used to sign the JWT token.
//
// Returns:
//   - string: the username included in the JWT claims.
//   - error: an error if the JWT token is not valid.
func ValidateJWTToken(tokenString string, signingKey []byte) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims["username"].(string), nil
	}
	return "", err
}
