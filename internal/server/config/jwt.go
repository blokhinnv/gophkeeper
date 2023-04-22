package config

import "time"

// jwtConfig is a part of the config which contains setting for the JWT tokens.
type jwtConfig struct {
	SigningKey     string        `env:"JWT_SIGNING_KEY"     envDefault:"practicum"`
	ExpireDuration time.Duration `env:"JWT_EXPIRE_DURATION" envDefault:"1h"`
}
