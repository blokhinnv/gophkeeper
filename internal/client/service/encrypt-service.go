package service

import (
	"encoding/json"
	"os"

	"github.com/blokhinnv/gophkeeper/pkg/encrypt"
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

	ciphertext, err := encrypt.EncryptBytes(data, key)
	if err != nil {
		return err
	}

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
	decoded, err := encrypt.DecryptBytes(ciphertext, key)
	if err != nil {
		return nil, err
	}

	resp := new(syncResponse)
	err = json.Unmarshal(decoded, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
