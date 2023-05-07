package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	srvErrors "github.com/blokhinnv/gophkeeper/internal/server/errors"
	"github.com/blokhinnv/gophkeeper/internal/server/models"
)

// SyncService defines the interface for syncing data.
type SyncService interface {
	// Sync syncs data from collections.
	Sync(token string, collections []models.CollectionName) (*syncResponse, error)
	Register(token, sockAddr string) (string, error)
	Unregister(token, sockAddr string) (string, error)
	// GetClient returns the service's client.
	GetClient() *resty.Client
}

// syncService implements the SyncService interface.
type syncService struct {
	client *resty.Client
}

// NewSyncService returns a new instance of SyncService.
func NewSyncService(baseURL string) SyncService {
	client := newConfiguredClient(baseURL)
	return &syncService{client: client}
}

// syncResponse defines the response from the SyncService Sync method.
type syncResponse struct {
	Text       []models.TextRecord
	Binary     []models.BinaryRecord
	Card       []models.CardRecord
	Credential []models.CredentialRecord
}

// Sync syncs data from collections.
func (s *syncService) Sync(
	token string,
	collectionNames []models.CollectionName,
) (*syncResponse, error) {
	r := &syncResponse{}
	for _, collectionName := range collectionNames {
		req := s.client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", fmt.Sprintf("Bearer: %v", token))
		switch collectionName {
		case models.TextCollection:
			req = req.SetResult(&r.Text)
		case models.BinaryCollection:
			req = req.SetResult(&r.Binary)
		case models.CardCollection:
			req = req.SetResult(&r.Card)
		case models.CredentialsCollection:
			req = req.SetResult(&r.Credential)
		default:
			return nil, fmt.Errorf("%w: %v", srvErrors.ErrUnknownCollection, collectionName)
		}
		resp, err := req.Get(fmt.Sprintf("/api/store/%v", collectionName))
		if err != nil {
			return nil, err
		}
		if resp.StatusCode() == http.StatusUnauthorized {
			return nil, srvErrors.ErrUnauthorized
		}
		if resp.StatusCode() >= http.StatusBadRequest {
			return nil, errors.New(resp.String())
		}
	}
	return r, nil
}

func (s *syncService) Register(token, sockAddr string) (string, error) {
	req := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer: %v", token)).
		SetBody(fmt.Sprintf(`{"socket_addr": "%v"}`, sockAddr))
	resp, err := req.Post("/api/sync/register")
	if err != nil {
		return "", err
	}
	if resp.StatusCode() >= http.StatusBadRequest {
		return "", errors.New(resp.String())
	}
	return resp.String(), nil
}

func (s *syncService) Unregister(token, sockAddr string) (string, error) {
	req := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer: %v", token)).
		SetBody(fmt.Sprintf(`{"socket_addr": "%v"}`, sockAddr))
	resp, err := req.Post("/api/sync/unregister")
	if err != nil {
		return "", err
	}
	if resp.StatusCode() >= http.StatusBadRequest {
		return "", errors.New(resp.String())
	}
	return resp.String(), nil
}

// GetClient returns the service's client.
func (s *syncService) GetClient() *resty.Client {
	return s.client
}
