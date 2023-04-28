package commands

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/blokhinnv/gophkeeper/internal/client/commands/auth"
	"github.com/blokhinnv/gophkeeper/internal/client/commands/crud"
	"github.com/blokhinnv/gophkeeper/internal/client/commands/shell"
	"github.com/blokhinnv/gophkeeper/internal/client/commands/sync"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "a client for a gophkeeper server",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(auth.AuthCmd)
	rootCmd.AddCommand(crud.CRUDCmd)
	rootCmd.AddCommand(shell.ShellCmd)
	rootCmd.AddCommand(sync.SyncCmd)
}
