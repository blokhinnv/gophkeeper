// Package main is an entry point for a client.
package main

import (
	"fmt"

	"github.com/blokhinnv/gophkeeper/internal/client"
)

// Global build variables
var (
	buildVersion string
	buildDate    string
)

// printBuildInfo prints a message containg build information.
func printBuildInfo() {
	coalesce := func(args ...string) string {
		for _, a := range args {
			if a != "" {
				return a
			}
		}
		return ""
	}
	buildVersion = coalesce(buildVersion, "N/A")
	fmt.Printf(
		"Build version: %s\nBuild date: %s\n",
		coalesce(buildVersion, "N/A"),
		coalesce(buildDate, "N/A"),
	)
}

func main() {
	printBuildInfo()
	client.RunClient()
}
