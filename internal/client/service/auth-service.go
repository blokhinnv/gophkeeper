// Package service provides different services for a client.
package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// AuthService is an interface that provides methods for authentication and registration.
type AuthService interface {
	// Auth authenticates a user with the given username and password and returns an authentication token if successful.
	Auth(username, password string) (string, error)
	// Register creates a new user with the given username and password.
	Register(username, password string) error
}

// authService is a concrete implementation of AuthService.
type authService struct {
	client *resty.Client
}

// NewAuthService creates a new instance of authService with the given baseURL and returns it as an AuthService.
func NewAuthService(baseURL string) AuthService {
	client := resty.New().SetBaseURL(baseURL)
	return &authService{client: client}
}

// authResult represents the JSON response for authentication and registration requests.
type authResult struct {
	Token string `json:"tok"`
	Error string `json:"error"`
}

// Auth authenticates a user with the given username and password and returns an authentication token if successful.
func (s *authService) Auth(username, password string) (string, error) {
	r := &authResult{}
	resp, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)).
		Put("/api/user/login")
	if err != nil {
		return "", err
	}
	if resp.StatusCode() >= http.StatusBadRequest {
		r.Error = resp.String()
	}
	r.Token = resp.String()
	if r.Error != "" {
		return "", errors.New(r.Error)
	}
	return r.Token, nil
}

// Register creates a new user with the given username and password.
func (s *authService) Register(username, password string) error {
	r := &authResult{}
	resp, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)).
		Put("/api/user/register")
	if err != nil {
		return err
	}
	if resp.StatusCode() >= http.StatusBadRequest {
		r.Error = resp.String()
	}
	r.Token = resp.String()
	if r.Error != "" {
		return errors.New(r.Error)
	}
	return nil
}
