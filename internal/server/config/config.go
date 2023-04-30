// Package config provides a configuration structure for a server application.
package config

import "github.com/caarlos0/env/v6"

// ServerConfig is a configuration struct for a server application that includes
// both database and JWT configurations.
type ServerConfig struct {
	dbConfig
	jwtConfig
	netConfig
}

// NewServerConfig creates a new ServerConfig object and populates its fields
// with values parsed from environment variables using the caarlos0/env/v6
// package. If parsing any of the environment variables fails, an error is
// returned.
func NewServerConfig() (*ServerConfig, error) {
	cfg := ServerConfig{}
	if err := env.Parse(&cfg.netConfig); err != nil {
		return nil, err
	}
	if err := env.Parse(&cfg.dbConfig); err != nil {
		return nil, err
	}
	if err := env.Parse(&cfg.jwtConfig); err != nil {
		return nil, err
	}
	return &cfg, nil
}
