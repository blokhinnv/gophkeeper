package service

import (
	"encoding/json"
	"os"

	"github.com/blokhinnv/gophkeeper/internal/client/models"
	"github.com/blokhinnv/gophkeeper/pkg/encrypt"
)

// EncryptService is an interface for encrypting and decrypting data to and from files.
type EncryptService interface {
	ToEncryptedFile(resp *models.SyncResponse, fileName, password string) error
	FromEncryptedFile(fileName, password string) (*models.SyncResponse, error)
}

// encryptService is the implementation of the EncryptService interface.
type encryptService struct {
}

// NewEncryptService creates a new instance of the EncryptService.
func NewEncryptService() EncryptService {
	return &encryptService{}
}

// ToEncryptedFile encrypts and writes models.SyncResponse data to a file.
// It uses AES encryption with the given password to encrypt the data.
// Returns an error if encryption or writing to file fails.
func (s *encryptService) ToEncryptedFile(resp *models.SyncResponse, fileName, key string) error {
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

// FromEncryptedFile reads and decrypts models.SyncResponse data from a file.
// It uses AES encryption with the given password to decrypt the data.
// Returns an error if decryption or reading from file fails.
func (s *encryptService) FromEncryptedFile(fileName, key string) (*models.SyncResponse, error) {
	ciphertext, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	decoded, err := encrypt.DecryptBytes(ciphertext, key)
	if err != nil {
		return nil, err
	}

	resp := new(models.SyncResponse)
	err = json.Unmarshal(decoded, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
