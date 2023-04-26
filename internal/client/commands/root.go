package commands

import (
	"os"

	"github.com/spf13/cobra"

	"gophkeeper/internal/client/commands/auth"
	"gophkeeper/internal/client/commands/item"
	"gophkeeper/internal/client/commands/shell"
	"gophkeeper/internal/client/commands/sync"
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
	rootCmd.AddCommand(item.ItemCmd)
	rootCmd.AddCommand(shell.ShellCmd)
	rootCmd.AddCommand(sync.SyncCmd)
}
