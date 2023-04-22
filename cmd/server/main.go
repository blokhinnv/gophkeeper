// Package main is an entry point for a server.
package main

import (
	"gophkeeper/internal/server"
	"gophkeeper/internal/server/config"
	"log"

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
