package service

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type StorageService interface {
	GetAll()
	Add(body, collection, token string) error
	Update(body, collection, token string) error
	Delete(body, collection, token string) error
}
type errResult struct {
	Error string `json:"error"`
}

type storageService struct {
	client *resty.Client
}

func NewStorageService(baseURL string) StorageService {
	client := resty.New().SetBaseURL(baseURL)
	return &storageService{client: client}
}

func (s *storageService) GetAll() {}

func (s *storageService) Add(body, collection, token string) error {
	r := &errResult{}
	_, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer: %v", token)).
		SetBody(body).
		SetError(r).
		Put(fmt.Sprintf("/api/store/%v", collection))
	if err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	return nil
}

func (s *storageService) Update(body, collection, token string) error {
	r := &errResult{}
	_, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer: %v", token)).
		SetBody(body).
		SetError(r).
		Post(fmt.Sprintf("/api/store/%v", collection))
	if err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	return nil
}

func (s *storageService) Delete(body, collection, token string) error {
	r := &errResult{}
	_, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer: %v", token)).
		SetBody(body).
		SetError(r).
		Delete(fmt.Sprintf("/api/store/%v", collection))
	if err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	return nil
}
