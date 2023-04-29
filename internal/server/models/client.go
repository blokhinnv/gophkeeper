package models

// Client represents a client connected to the server.
type Client struct {
	Username   string `json:"-"`
	SocketAddr string `json:"socket_addr"`
}
