package config

import "time"

// jwtConfig is a part of the config which contains setting for the JWT tokens.
type jwtConfig struct {
	SigningKey     string        `env:"GOPHKEEPER_JWT_SIGNING_KEY"     envDefault:"practicum"`
	ExpireDuration time.Duration `env:"GOPHKEEPER_JWT_EXPIRE_DURATION" envDefault:"1h"`
}
