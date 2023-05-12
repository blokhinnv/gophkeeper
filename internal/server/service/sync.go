package service

import (
	"net"
	"sync"

	"golang.org/x/exp/slices"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"github.com/blokhinnv/gophkeeper/internal/server/models"
	slicesUtils "github.com/blokhinnv/gophkeeper/pkg/slices"
)

// SyncService is an interface for syncing clients.
type SyncService interface {
	Register(client *models.Client)
	Unregister(client *models.Client)
	Signal(client *models.Client)
}

// syncService implements the SyncService interface.
type syncService struct {
	client  *resty.Client
	clients map[string][]string // Addr: Username
	mu      sync.Mutex
}

// NewSyncService creates a new SyncService instance.
func NewSyncService() SyncService {
	s := &syncService{
		client:  resty.New(),
		clients: make(map[string][]string),
	}
	return s
}

// Register registers a new client with the server.
func (s *syncService) Register(client *models.Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	addrs := s.clients[client.Username]
	if !slices.Contains(addrs, client.SocketAddr) {
		s.clients[client.Username] = append(addrs, client.SocketAddr)
	}
	logrus.Infof("Registered %v %v", client.Username, client.SocketAddr)
}

// Unregister unregisters an existing client from the server.
func (s *syncService) Unregister(client *models.Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	oldAddrs, ok := s.clients[client.Username]
	if !ok {
		return
	}
	s.clients[client.Username] = slicesUtils.Remove(oldAddrs, client.SocketAddr)
	logrus.Infof("Unregistered %v %v", client.Username, client.SocketAddr)
}

// Signal sends a signal to all the user's registered clients.
func (s *syncService) Signal(client *models.Client) {
	for _, addr := range s.clients[client.Username] {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			s.clients[client.Username] = slicesUtils.Remove(s.clients[client.Username], addr)
			logrus.Infof("unable to reach %v; unregistered", addr)
			continue
		}
		_, err = conn.Write(nil)
		if err != nil {
			s.clients[client.Username] = slicesUtils.Remove(s.clients[client.Username], addr)
			logrus.Infof("unable to reach %v; unregistered", addr)
		}
	}
}

// getClientAddresses returns addresses for the username provided.
func (m *syncService) getClientAddresses(username string) []string {
	return m.clients[username]
}
