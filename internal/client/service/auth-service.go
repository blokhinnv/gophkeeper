package service

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type AuthService interface {
	Auth(username, password string) (string, error)
	Register(username, password string) error
}

type authService struct {
	client *resty.Client
}

func NewAuthService(baseURL string) AuthService {
	client := resty.New().SetBaseURL(baseURL)
	return &authService{client: client}
}

type authResult struct {
	Token string `json:"tok"`
	Error string `json:"error"`
}

func (s *authService) Auth(username, password string) (string, error) {
	r := &authResult{}
	_, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)).
		SetResult(r).
		SetError(r).
		Put("/api/user/login")
	if err != nil {
		return "", err
	}
	if r.Error != "" {
		return "", errors.New(r.Error)
	}
	return r.Token, nil
}

func (s *authService) Register(username, password string) error {
	r := &authResult{}
	_, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)).
		SetResult(r).
		SetError(r).
		Put("/api/user/register")
	if err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	return nil
}
