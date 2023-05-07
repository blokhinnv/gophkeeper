package service

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestNewConfiguredClient(t *testing.T) {
	t.Run("http", func(t *testing.T) {
		baseURL := "http://example.com"
		client := newConfiguredClient(baseURL)
		assert.IsType(t, &resty.Client{}, client)
		assert.Equal(t, baseURL, client.HostURL)
	})
	t.Run("https", func(t *testing.T) {
		baseURL := "https://example.com"
		client := newConfiguredClient(baseURL)
		assert.IsType(t, &resty.Client{}, client)
		assert.Equal(t, baseURL, client.HostURL)
		assert.NotNil(t, client.GetClient().Transport.(*http.Transport).TLSClientConfig)
		assert.True(
			t,
			client.GetClient().Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify,
		)
	})
}

func TestAuthService_Auth(t *testing.T) {
	baseURL := "https://example.com"
	service := NewAuthService(baseURL)
	assert.NotNil(t, service)
	client := service.GetClient()
	assert.Equal(t, baseURL, client.HostURL)

	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()
	t.Run("ok", func(t *testing.T) {
		httpmock.Reset()

		httpmock.RegisterResponder(
			http.MethodPut,
			fmt.Sprintf("%v/api/user/login", baseURL),
			httpmock.NewStringResponder(200, "some token..."),
		)

		resp, err := service.Auth("testuser", "testpassword")
		assert.Equal(t, "some token...", resp)
		assert.Nil(t, err)
	})
	t.Run("fail", func(t *testing.T) {
		httpmock.Reset()

		httpmock.RegisterResponder(
			http.MethodPut,
			fmt.Sprintf("%v/api/user/login", baseURL),
			httpmock.NewStringResponder(400, "some error..."),
		)

		resp, err := service.Auth("testuser", "testpassword")
		assert.Equal(t, "some error...", err.Error())
		assert.NotNil(t, err)
		assert.Equal(t, "", resp)
	})
}

func TestAuthService_Register(t *testing.T) {
	baseURL := "https://example.com"
	service := NewAuthService(baseURL)
	assert.NotNil(t, service)
	client := service.GetClient()
	assert.Equal(t, baseURL, client.HostURL)

	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()
	t.Run("ok", func(t *testing.T) {
		httpmock.Reset()

		httpmock.RegisterResponder(
			http.MethodPut,
			fmt.Sprintf("%v/api/user/register", baseURL),
			httpmock.NewStringResponder(200, "ok"),
		)

		err := service.Register("testuser", "testpassword")
		assert.Nil(t, err)
	})
	t.Run("fail", func(t *testing.T) {
		httpmock.Reset()

		httpmock.RegisterResponder(
			http.MethodPut,
			fmt.Sprintf("%v/api/user/register", baseURL),
			httpmock.NewStringResponder(400, "some error..."),
		)

		err := service.Register("testuser", "testpassword")
		assert.NotNil(t, err)
	})
}
