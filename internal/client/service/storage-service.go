package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	clientErr "github.com/blokhinnv/gophkeeper/internal/client/errors"
	clientModels "github.com/blokhinnv/gophkeeper/internal/client/models"
	srvrModels "github.com/blokhinnv/gophkeeper/internal/server/models"
)

// StorageService defines the interface for managing data storage.
type StorageService interface {
	// GetAll retrieves all data from a specific collection.
	GetAll(collectionName srvrModels.CollectionName, data *clientModels.SyncResponse) any
	// Add adds a new item to a specific collection.
	Add(body string, collectionName srvrModels.CollectionName, token string) (string, error)
	// Update updates an existing item in a specific collection.
	Update(body string, collectionName srvrModels.CollectionName, token string) (string, error)
	// Delete removes an existing item from a specific collection.
	Delete(body string, collectionName srvrModels.CollectionName, token string) (string, error)
	// GetClient returns the service's client.
	GetClient() *resty.Client
}

// storageService is an implementation of the StorageService interface.
type storageService struct {
	client *resty.Client
}

// NewStorageService returns a new instance of StorageService.
func NewStorageService(baseURL string) StorageService {
	client := newConfiguredClient(baseURL)
	return &storageService{client: client}
}

// GetAll retrieves all data from a specific collection.
func (s *storageService) GetAll(
	collectionName srvrModels.CollectionName,
	data *clientModels.SyncResponse,
) any {
	switch collectionName {
	case srvrModels.TextCollection:
		return data.Text
	case srvrModels.BinaryCollection:
		return data.Binary
	case srvrModels.CardCollection:
		return data.Card
	case srvrModels.CredentialsCollection:
		return data.Credential
	default:
		return nil
	}
}

// Add adds a new item to a specific collection.
func (s *storageService) Add(
	body string,
	collectionName srvrModels.CollectionName,
	token string,
) (string, error) {
	resp, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer: %v", token)).
		SetBody(body).
		Put(fmt.Sprintf("/api/store/%v", collectionName))
	if err != nil {
		return "", fmt.Errorf("%w: %v", clientErr.ErrServerUnavailable, err)
	}
	if resp.StatusCode() >= http.StatusBadRequest {
		return "", errors.New(resp.String())
	}
	return resp.String(), nil
}

// Update updates an existing item in a specific collection.
func (s *storageService) Update(
	body string,
	collectionName srvrModels.CollectionName,
	token string,
) (string, error) {
	resp, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer: %v", token)).
		SetBody(body).
		Post(fmt.Sprintf("/api/store/%v", collectionName))
	if err != nil {
		return "", fmt.Errorf("%w: %v", clientErr.ErrServerUnavailable, err)
	}
	if resp.StatusCode() >= http.StatusBadRequest {
		return "", errors.New(resp.String())
	}
	return resp.String(), nil
}

// Delete removes an existing item from a specific collection.
func (s *storageService) Delete(
	body string,
	collectionName srvrModels.CollectionName,
	token string,
) (string, error) {
	resp, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer: %v", token)).
		SetBody(body).
		Delete(fmt.Sprintf("/api/store/%v", collectionName))
	if err != nil {
		return "", fmt.Errorf("%w: %v", clientErr.ErrServerUnavailable, err)
	}
	if resp.StatusCode() >= http.StatusBadRequest {
		return "", errors.New(resp.String())
	}
	return resp.String(), nil
}

// GetClient returns the service's client.
func (s *storageService) GetClient() *resty.Client {
	return s.client
}
