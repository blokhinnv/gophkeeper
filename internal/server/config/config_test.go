package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewServerConfig(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("GOPHKEEPER_DB_URI", "mongodb://test-db:27017")
	os.Setenv("GOPHKEEPER_DB_NAME", "test-db-name")
	os.Setenv("GOPHKEEPER_DB_ENCRYPTION_KEY", "test-encryption-key")
	os.Setenv("JWT_SIGNING_KEY", "test-signing-key")
	os.Setenv("JWT_EXPIRE_DURATION", "2h")
	os.Setenv("GOPHKEEPER_SERVER_PORT", "8888")
	os.Setenv("GOPHKEEPER_USE_HTTPS", "false")
	os.Setenv("GOPHKEEPER_CERT_FILE", "test-cert-file")
	os.Setenv("GOPHKEEPER_KEY_FILE", "test-key-file")

	// Cleanup environment variables after the test
	defer func() {
		os.Unsetenv("GOPHKEEPER_DB_URI")
		os.Unsetenv("GOPHKEEPER_DB_NAME")
		os.Unsetenv("GOPHKEEPER_DB_ENCRYPTION_KEY")
		os.Unsetenv("JWT_SIGNING_KEY")
		os.Unsetenv("JWT_EXPIRE_DURATION")
		os.Unsetenv("GOPHKEEPER_SERVER_PORT")
		os.Unsetenv("GOPHKEEPER_USE_HTTPS")
		os.Unsetenv("GOPHKEEPER_CERT_FILE")
		os.Unsetenv("GOPHKEEPER_KEY_FILE")
	}()

	expected := &ServerConfig{
		dbConfig: dbConfig{
			MongoURI:      "mongodb://test-db:27017",
			DBName:        "test-db-name",
			EncryptionKey: "test-encryption-key",
		},
		jwtConfig: jwtConfig{
			SigningKey:     "test-signing-key",
			ExpireDuration: 2 * time.Hour,
		},
		netConfig: netConfig{
			Port:     "8888",
			UseHTTPS: false,
			CertFile: "test-cert-file",
			KeyFile:  "test-key-file",
		},
	}

	// Call NewServerConfig to get the actual value
	actual, err := NewServerConfig()
	require.NoError(t, err)

	// Compare the actual and expected values
	require.Equal(t, expected, actual)
}
