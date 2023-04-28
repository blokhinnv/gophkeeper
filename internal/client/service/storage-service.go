package service

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"

	"gophkeeper/internal/server/models"
)

// StorageService defines the interface for managing data storage.
type StorageService interface {
	// GetAll retrieves all data from a specific collection.
	GetAll(collectionName models.Collection, data *syncResponse) any
	// Add adds a new item to a specific collection.
	Add(body string, collectionName models.Collection, token string) error
	// Update updates an existing item in a specific collection.
	Update(body string, collectionName models.Collection, token string) error
	// Delete removes an existing item from a specific collection.
	Delete(body string, collectionName models.Collection, token string) error
}

// errResult represents the response body for API requests that may return an error.
type errResult struct {
	Error string `json:"error"`
}

// storageService is an implementation of the StorageService interface.
type storageService struct {
	client *resty.Client
}

// NewStorageService returns a new instance of StorageService.
func NewStorageService(baseURL string) StorageService {
	client := resty.New().SetBaseURL(baseURL)
	return &storageService{client: client}
}

// GetAll retrieves all data from a specific collection.
func (s *storageService) GetAll(collectionName models.Collection, data *syncResponse) any {
	switch collectionName {
	case models.TextCollection:
		return data.Text
	case models.BinaryCollection:
		return data.Binary
	case models.CardCollection:
		return data.Card
	case models.CredentialsCollection:
		return data.Credential
	default:
		return nil
	}
}

// Add adds a new item to a specific collection.
func (s *storageService) Add(
	body string,
	collectionName models.Collection,
	token string,
) error {
	r := &errResult{}
	_, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer: %v", token)).
		SetBody(body).
		SetError(r).
		Put(fmt.Sprintf("/api/store/%v", collectionName))
	if err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	return nil
}

// Update updates an existing item in a specific collection.
func (s *storageService) Update(
	body string,
	collectionName models.Collection,
	token string,
) error {
	r := &errResult{}
	_, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer: %v", token)).
		SetBody(body).
		SetError(r).
		Post(fmt.Sprintf("/api/store/%v", collectionName))
	if err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	return nil
}

// Delete removes an existing item from a specific collection.
func (s *storageService) Delete(
	body string,
	collectionName models.Collection,
	token string,
) error {
	r := &errResult{}
	_, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer: %v", token)).
		SetBody(body).
		SetError(r).
		Delete(fmt.Sprintf("/api/store/%v", collectionName))
	if err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	return nil
}
