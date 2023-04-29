// Package server provides functionality for running a client for a gophkeeper server.
package client

import "github.com/blokhinnv/gophkeeper/internal/client/commands"

func RunClient() {
	commands.Execute()
}
