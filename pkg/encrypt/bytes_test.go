package encrypt

import (
	"bytes"
	"crypto/aes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenToAESKey(t *testing.T) {
	// Call the tokenToAESKey function
	aesKey := tokenToAESKey("test-key")
	assert.Equal(t, 16, len(aesKey))
}

func TestEncryptBytes(t *testing.T) {
	t.Run("empty key", func(t *testing.T) {
		data := []byte("secret message")
		key := ""
		_, err := EncryptBytes(data, key)
		if err == nil {
			t.Fatalf("expected an error, but got none")
		}
	})
	t.Run("empty data", func(t *testing.T) {
		data := []byte{}
		key := "mysecretkey"
		_, err := EncryptBytes(data, key)
		if err == nil {
			t.Fatalf("expected an error, but got none")
		}
	})
	t.Run("random", func(t *testing.T) {
		data1 := []byte("secret message")
		data2 := []byte("another secret message")
		key := "mysecretkey"

		encrypted1, err := EncryptBytes(data1, key)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		encrypted2, err := EncryptBytes(data2, key)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// make sure the encrypted data is not the same for two different inputs
		if bytes.Equal(encrypted1, encrypted2) {
			t.Errorf("encrypted data should not be the same for two different inputs")
		}
	})
}

func TestDecryptBytes(t *testing.T) {
	key := "this-is-a-secret-key"
	message := "hello world"
	encrypted, err := EncryptBytes([]byte(message), key)
	if err != nil {
		t.Fatalf("error encrypting message: %v", err)
	}
	t.Run("correct key", func(t *testing.T) {
		// test correct decryption
		decrypted, err := DecryptBytes(encrypted, key)
		if err != nil {
			t.Fatalf("error decrypting message: %v", err)
		}
		if !bytes.Equal(decrypted, []byte(message)) {
			t.Errorf(
				"decrypted message does not match original message: %v != %v",
				decrypted,
				[]byte(message),
			)
		}
	})
	t.Run("incorrect_key", func(t *testing.T) {
		// test incorrect key
		decrypted, _ := DecryptBytes(encrypted, "wrong-key")
		if bytes.Equal(decrypted, []byte(message)) {
			t.Errorf(
				"decrypted message should not match original message: %v != %v",
				decrypted,
				[]byte(message),
			)
		}
	})
	t.Run("short_cipher", func(t *testing.T) {
		// test short ciphertext
		shortCiphertext := make([]byte, aes.BlockSize-1)
		_, err = DecryptBytes(shortCiphertext, key)
		if err == nil {
			t.Errorf("expected error with short ciphertext, but got nil")
		}
	})
}
