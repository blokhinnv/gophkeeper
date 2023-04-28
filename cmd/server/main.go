// Package main is an entry point for a server.
package main

import (
	"log"

	"github.com/blokhinnv/gophkeeper/internal/server"
	"github.com/blokhinnv/gophkeeper/internal/server/config"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	cfg, err := config.NewServerConfig()
	if err != nil {
		log.Fatal(err)
	}
	server.RunServer(cfg)
}
