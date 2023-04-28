package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

func main() {
	client := resty.New()
	r, err := client.SetBaseURL("http://localhost:8081").R().Get("/")
	logrus.Info("Polling ", r, err)
}
