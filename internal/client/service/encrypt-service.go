package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

// EncryptService is an interface for encrypting and decrypting data to and from files.
type EncryptService interface {
	ToEncryptedFile(resp *syncResponse, fileName, password string) error
	FromEncryptedFile(fileName, password string) (*syncResponse, error)
}

// encryptService is the implementation of the EncryptService interface.
type encryptService struct {
}

// NewEncryptService creates a new instance of the EncryptService.
func NewEncryptService() EncryptService {
	return &encryptService{}
}

// ToEncryptedFile encrypts and writes syncResponse data to a file.
// It uses AES encryption with the given password to encrypt the data.
// Returns an error if encryption or writing to file fails.
func (s *encryptService) ToEncryptedFile(resp *syncResponse, fileName, key string) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	aesKey, err := s.tokenToAESKey(key)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return err
	}
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}
	mode := cipher.NewCFBEncrypter(block, iv)
	mode.XORKeyStream(ciphertext[aes.BlockSize:], data)

	err = os.WriteFile(fileName, ciphertext, 0644)
	if err != nil {
		return err
	}
	return nil
}

// FromEncryptedFile reads and decrypts syncResponse data from a file.
// It uses AES encryption with the given password to decrypt the data.
// Returns an error if decryption or reading from file fails.
func (s *encryptService) FromEncryptedFile(fileName, key string) (*syncResponse, error) {
	ciphertext, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	aesKey, err := s.tokenToAESKey(key)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	resp := new(syncResponse)
	err = json.Unmarshal(ciphertext, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// tokenToAESKey derives an AES key using PBKDF2 with SHA256 as the hash function.
func (s *encryptService) tokenToAESKey(key string) ([]byte, error) {
	aesKey := pbkdf2.Key([]byte(key), nil, 1000, 16, sha256.New)
	return aesKey, nil
}
