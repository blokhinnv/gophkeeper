// Package config provides a configuration structure for a server application.
package config

import "github.com/caarlos0/env/v6"

// ServerConfig is a configuration struct for a server application that includes
// both database and JWT configurations.
type ServerConfig struct {
	dbConfig
	jwtConfig
}

// NewServerConfig creates a new ServerConfig object and populates its fields
// with values parsed from environment variables using the caarlos0/env/v6
// package. If parsing any of the environment variables fails, an error is
// returned.
//
// Returns:
//   - *ServerConfig: a new ServerConfig object with fields populated with
//     environment variable values.
//   - error: an error if any of the environment variables fail to parse.
func NewServerConfig() (*ServerConfig, error) {
	cfg := ServerConfig{}
	if err := env.Parse(&cfg.dbConfig); err != nil {
		return nil, err
	}
	if err := env.Parse(&cfg.jwtConfig); err != nil {
		return nil, err
	}
	return &cfg, nil
}
