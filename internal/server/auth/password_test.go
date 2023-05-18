package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestVerifyPassword(t *testing.T) {
	password := "myPassword123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("unexpected error while generating bcrypt hash: %v", err)
	}

	// Test case: matching password and hash
	match, err := VerifyPassword(password, string(hash))
	if err != nil {
		t.Fatalf("unexpected error while verifying password: %v", err)
	}
	if !match {
		t.Errorf("password should match hash, but VerifyPassword returned false")
	}

	// Test case: non-matching password and hash
	match, err = VerifyPassword("wrongPassword", string(hash))
	if err != nil {
		t.Fatalf("unexpected error while verifying password: %v", err)
	}
	if match {
		t.Errorf("password should not match hash, but VerifyPassword returned true")
	}

	// Test case: invalid hash format
	match, err = VerifyPassword(password, "invalidHashFormat")
	if err == nil {
		t.Errorf(
			"expected error while verifying invalid hash format, but VerifyPassword returned no error",
		)
	}
	if match {
		t.Errorf("password should not match invalid hash format, but VerifyPassword returned true")
	}
}
