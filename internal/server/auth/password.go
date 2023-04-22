package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// VerifyPassword verifies the provided plaintext password against a bcrypt hash. The function
// returns a boolean indicating whether the password matches the hash or not. If the password
// and hash do not match, the function returns `false`. If an error occurs during the comparison,
// the function returns `false` and the error.
//
// Parameters:
//   - password: the plaintext password to verify.
//   - hash: the bcrypt hash to verify the password against.
//
// Returns:
//   - bool: `true` if the password matches the hash, `false` otherwise.
//   - error: an error if an error occurs during the comparison.
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
