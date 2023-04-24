package service

import (
	"gophkeeper/internal/server/models"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// SyncService is an interface for syncing clients.
type SyncService interface {
	Register(client *models.Client)
	Unregister(client *models.Client)
}

// syncService implements the SyncService interface.
type syncService struct {
	client  *resty.Client
	clients map[string]string // Addr: Username
	mu      sync.Mutex
}

// NewSyncService creates a new SyncService instance.
func NewSyncService() SyncService {
	s := &syncService{
		client:  resty.New(),
		clients: make(map[string]string),
	}
	go s.poll()
	return s
}

// Register registers a new client with the server.
func (s *syncService) Register(client *models.Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	logrus.Infof("Registered %v", client.Addr)
	s.clients[client.Addr] = client.Username
}

// Unregister unregisters an existing client from the server.
func (s *syncService) Unregister(client *models.Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	logrus.Infof("Unregistered %v", client.Addr)
	delete(s.clients, client.Addr)
}

// poll continuously polls registered clients to check if they are still available.
// It is run as a goroutine when NewSyncService is called.
func (s *syncService) poll() {
	ticker := time.NewTicker(500 * time.Millisecond)
	for range ticker.C {
		for addr := range s.clients {
			_, err := s.client.SetBaseURL(addr).R().Get("/")
			logrus.Info("Polling ", err)
			if err != nil {
				logrus.Infof("Unregistered %v", addr)
				delete(s.clients, addr)
			}
		}
	}
}
