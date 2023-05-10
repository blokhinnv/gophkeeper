// Package encrypt provides functions for encrypting the data.
package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/xdg-go/pbkdf2"
)

// tokenToAESKey derives an AES key using PBKDF2 with SHA256 as the hash function.
func tokenToAESKey(key string) []byte {
	aesKey := pbkdf2.Key([]byte(key), nil, 1000, 16, sha256.New)
	return aesKey
}

func EncryptBytes(data []byte, key string) ([]byte, error) {
	if key == "" {
		return nil, fmt.Errorf("empty key")
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}
	aesKey := tokenToAESKey(key)
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	mode := cipher.NewCFBEncrypter(block, iv)
	mode.XORKeyStream(ciphertext[aes.BlockSize:], data)
	return ciphertext, nil
}

func DecryptBytes(encryptedData []byte, key string) ([]byte, error) {
	aesKey := tokenToAESKey(key)
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	ciphertext := []byte(encryptedData)
	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext, nil
}
