package config

// netConfig is a part of the config which contains setting for network.
type netConfig struct {
	Port     string `env:"GOPHKEEPER_SERVER_PORT" envDefault:"8080"`
	UseHTTPS bool   `env:"GOPHKEEPER_USE_HTTPS"   envDefault:"true"`
	CertFile string `env:"GOPHKEEPER_CERT_FILE"`
	KeyFile  string `env:"GOPHKEEPER_KEY_FILE"`
}
