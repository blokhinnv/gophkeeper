package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// VerifyPassword verifies the provided plaintext password against a bcrypt hash. The function
// returns a boolean indicating whether the password matches the hash or not. If the password
// and hash do not match, the function returns `false`. If an error occurs during the comparison,
// the function returns `false` and the error.
func VerifyPassword(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
