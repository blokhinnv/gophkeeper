package encrypt

import (
	"reflect"
	"testing"
)

func TestEncryptMap(t *testing.T) {
	key := "secretkey"
	t.Run("ok", func(t *testing.T) {
		data := map[string]any{
			"name":   "Alice",
			"age":    30,
			"gender": "female",
		}

		result, err := EncryptMap(data, key)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if data["name"] == result["name"] {
			t.Errorf("Result name %v should not match %v", data["name"], result["name"])
		}
		if data["gender"] == result["gender"] {
			t.Errorf("Result gender %v should not match %v", data["gender"], result["gender"])
		}
		if result["age"] != nil {
			t.Errorf("Result age %v should be nil", data["age"])
		}
	})
	t.Run("empty", func(t *testing.T) {
		data := make(map[string]any)
		expected := make(map[string]any)
		result, err := EncryptMap(data, key)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Result %v does not match expected %v", result, expected)
		}
	})
}

func TestDecryptMap(t *testing.T) {
	key := "test-key"

	data := map[string]any{
		"foo": "foo",
		"bar": "bar",
		"baz": "baz",
	}
	encryptedData, err := EncryptMap(data, key)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	decryptedData, err := DecryptMap(encryptedData, key)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(decryptedData) != len(encryptedData) {
		t.Errorf(
			"Expected decrypted data length of %v, but got %v",
			len(encryptedData),
			len(decryptedData),
		)
	}

	for k, v := range data {
		decryptedValue, ok := decryptedData[k]
		if !ok {
			t.Errorf("Expected decrypted data to contain key %v", k)
		}

		expectedValue, ok := v.(string)
		if !ok {
			t.Errorf("Unexpected error: %v", err)
		}

		if decryptedValue != expectedValue {
			t.Errorf(
				"Expected value of %v for key %v, but got %v",
				expectedValue,
				k,
				decryptedValue,
			)
		}
	}
}
