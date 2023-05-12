// Package main is an entry point for a server.
package main

import (
	"github.com/joho/godotenv"

	"github.com/blokhinnv/gophkeeper/internal/server"
	"github.com/blokhinnv/gophkeeper/internal/server/config"
	"github.com/blokhinnv/gophkeeper/pkg/log"
)

func main() {
	godotenv.Load(".env")
	cfg, err := config.NewServerConfig()
	if err != nil {
		log.Fatal(err)
	}
	server.RunServer(cfg)
}
