package models

import (
	"fmt"
	"testing"

	"github.com/blokhinnv/gophkeeper/internal/server/errors"
)

func TestNewCollectionName(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// Test a valid collection name
		name, err := NewCollectionName("text")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if name != TextCollection {
			t.Errorf("expected name to be %v, but got %v", TextCollection, name)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		// Test an invalid collection name
		invalidName := "invalid_collection_name"
		_, err := NewCollectionName(invalidName)
		if err == nil {
			t.Errorf("expected error, but got none")
		}
		expectedErr := fmt.Sprintf("%v: %v", errors.ErrUnknownCollection, invalidName)
		if err.Error() != expectedErr {
			t.Errorf("expected error message to be %v, but got %v", expectedErr, err)
		}
	})
}
