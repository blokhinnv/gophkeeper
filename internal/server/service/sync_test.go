package service

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/blokhinnv/gophkeeper/internal/server/models"
)

func TestNewSyncService(t *testing.T) {
	s := NewSyncService()
	assert.NotNil(t, s, "NewSyncService should return a non-nil pointer.")
}

func TestSyncService_Register(t *testing.T) {
	// Create a mock sync service instance.
	m := syncService{
		clients: make(map[string][]string),
	}

	// Test case 1: Registering a new client should add their address to the clients map.
	client := &models.Client{
		Username:   "testuser",
		SocketAddr: "127.0.0.1:1234",
	}
	m.Register(client)
	assert.Contains(
		t,
		m.clients["testuser"],
		"127.0.0.1:1234",
		"Register should add a new client address to the clients map.",
	)

	// Test case 2: Registering a client with an existing username should add their address to the existing slice.
	client2 := &models.Client{
		Username:   "testuser",
		SocketAddr: "127.0.0.1:5678",
	}
	m.Register(client2)
	assert.Equal(
		t,
		[]string{"127.0.0.1:1234", "127.0.0.1:5678"},
		m.getClientAddresses("testuser"),
		"Register should add a new client address to the existing slice in the clients map.",
	)

	// Test case 3: Registering a client with an existing username and address should not add a new entry to the clients map.
	client3 := &models.Client{
		Username:   "testuser",
		SocketAddr: "127.0.0.1:5678",
	}
	m.Register(client3)
	assert.Equal(
		t,
		[]string{"127.0.0.1:1234", "127.0.0.1:5678"},
		m.getClientAddresses("testuser"),
		"Register should not add a new entry to the clients map if the username and address already exist.",
	)
}

func TestSyncService_Unregister(t *testing.T) {
	// Create a mock sync service instance with a test client.
	m := &syncService{
		clients: map[string][]string{
			"testuser": {"127.0.0.1:1234", "127.0.0.1:5678"},
		},
	}

	// Test case 1: Unregistering an existing client should remove their address from the clients map.
	client := &models.Client{
		Username:   "testuser",
		SocketAddr: "127.0.0.1:1234",
	}
	m.Unregister(client)
	assert.NotContains(
		t,
		m.clients["testuser"],
		"127.0.0.1:1234",
		"Unregister should remove an existing client address from the clients map.",
	)

	// Test case 2: Unregistering a client with an unknown username should not modify the clients map.
	client2 := &models.Client{
		Username:   "unknownuser",
		SocketAddr: "127.0.0.1:5678",
	}
	m.Unregister(client2)
	assert.Equal(
		t,
		[]string{"127.0.0.1:5678"},
		m.getClientAddresses("testuser"),
		"Unregister should not modify the clients map if the username is unknown.",
	)

	// Test case 3: Unregistering a client with an unknown address for a known username should not modify the clients map.
	client3 := &models.Client{
		Username:   "testuser",
		SocketAddr: "127.0.0.1:9999",
	}
	m.Unregister(client3)
	assert.Equal(
		t,
		[]string{"127.0.0.1:5678"},
		m.getClientAddresses("testuser"),
		"Unregister should not modify the clients map if the address is unknown for a known username.",
	)
}

func TestSyncService_Signal(t *testing.T) {
	// Create a mock sync service instance with a test client.
	m := &syncService{
		clients: map[string][]string{
			"testuser": {"127.0.0.1:1234", "127.0.0.1:5678"},
		},
	}

	// Test case 1: Sending a signal to all registered clients should successfully write to the connection.
	client := &models.Client{
		Username:   "testuser",
		SocketAddr: "127.0.0.1:1234",
	}
	m.Register(client)
	// Create a mock TCP server that listens on the test client's socket address
	listener, err := net.Listen("tcp", client.SocketAddr)
	require.NoError(t, err)
	go func() {
		for {
			// Accept incoming connections
			conn, err := listener.Accept()
			if err != nil {
				return
			}

			// Close the connection immediately
			conn.Close()
		}
	}()
	// Send a signal to the test client
	m.Signal(client)

	// Verify that the test client's socket address is still registered
	assert.Equal(t, []string{client.SocketAddr}, m.getClientAddresses(client.Username))

	// Close the mock TCP server
	require.NoError(t, listener.Close())
}
