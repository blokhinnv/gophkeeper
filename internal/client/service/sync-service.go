package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	clientModels "github.com/blokhinnv/gophkeeper/internal/client/models"
	srvErrors "github.com/blokhinnv/gophkeeper/internal/server/errors"
	srvrModels "github.com/blokhinnv/gophkeeper/internal/server/models"
)

// SyncService defines the interface for syncing data.
type SyncService interface {
	// Sync syncs data from collections.
	Sync(token string, collections []srvrModels.CollectionName) (*clientModels.SyncResponse, error)
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

// Sync syncs data from collections.
func (s *syncService) Sync(
	token string,
	collectionNames []srvrModels.CollectionName,
) (*clientModels.SyncResponse, error) {
	r := &clientModels.SyncResponse{}
	for _, collectionName := range collectionNames {
		req := s.client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", fmt.Sprintf("Bearer: %v", token))
		switch collectionName {
		case srvrModels.TextCollection:
			req = req.SetResult(&r.Text)
		case srvrModels.BinaryCollection:
			req = req.SetResult(&r.Binary)
		case srvrModels.CardCollection:
			req = req.SetResult(&r.Card)
		case srvrModels.CredentialsCollection:
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

// Register registers a socket address with the synchronization service.
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

// Unregister removes a registered socket address from the synchronization service.
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
