package service

import (
	"net/http"
	"time"
)

type Service interface {
	ParseEnv() map[string]string
}

type service struct {
	client *http.Client
}

func NewService() (Service, error) {
	client := &http.Client{
		Timeout: time.Minute,
	}

	svc := &service{
		client: client,
	}
	return svc, nil
}
