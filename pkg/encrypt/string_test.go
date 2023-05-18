package encrypt

import (
	"crypto/aes"
	"testing"
)

func TestEncryptString(t *testing.T) {
	t.Run("empty key", func(t *testing.T) {
		data := "secret message"
		key := ""
		_, err := EncryptString(data, key)
		if err == nil {
			t.Fatalf("expected an error, but got none")
		}
	})
	t.Run("empty data", func(t *testing.T) {
		data := ""
		key := "mysecretkey"
		_, err := EncryptString(data, key)
		if err == nil {
			t.Fatalf("expected an error, but got none")
		}
	})
	t.Run("random", func(t *testing.T) {
		data1 := "secret message"
		data2 := "another secret message"
		key := "mysecretkey"

		encrypted1, err := EncryptString(data1, key)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		encrypted2, err := EncryptString(data2, key)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// make sure the encrypted data is not the same for two different inputs
		if encrypted1 == encrypted2 {
			t.Errorf("encrypted data should not be the same for two different inputs")
		}
	})
}

func TestDecryptString(t *testing.T) {
	key := "this-is-a-secret-key"
	message := "hello world"
	encrypted, err := EncryptString(message, key)
	if err != nil {
		t.Fatalf("error encrypting message: %v", err)
	}
	t.Run("correct key", func(t *testing.T) {
		// test correct decryption
		decrypted, err := DecryptString(encrypted, key)
		if err != nil {
			t.Fatalf("error decrypting message: %v", err)
		}
		if decrypted != message {
			t.Errorf(
				"decrypted message does not match original message: %v != %v",
				decrypted,
				[]byte(message),
			)
		}
	})
	t.Run("incorrect_key", func(t *testing.T) {
		// test incorrect key
		decrypted, _ := DecryptString(encrypted, "wrong-key")
		if decrypted == message {
			t.Errorf(
				"decrypted message should not match original message: %v != %v",
				decrypted,
				message,
			)
		}
	})
	t.Run("short_cipher", func(t *testing.T) {
		// test short ciphertext
		shortCiphertext := make([]byte, aes.BlockSize-1)
		_, err = DecryptString(string(shortCiphertext), key)
		if err == nil {
			t.Errorf("expected error with short ciphertext, but got nil")
		}
	})
}
