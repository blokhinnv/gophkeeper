// Package shell provides implementation of an interactive console client.
package shell

import (
	"github.com/spf13/cobra"
)

var (
	// shellCtrl is a controller to interact with different services.
	shellCtrl ShellController
	// ShellCmd represents the shell command.
	ShellCmd = &cobra.Command{
		Use:   "shell",
		Short: "Runs the shell with a persistent menu.",
		Long: `Runs the shell and presents a persistent menu for the user to interact with.
The user is presented with a menu that persists throughout the duration of the shell, a
llowing them to navigate through different features of the shell.`,
		Run: func(cmd *cobra.Command, args []string) {
			for {
				shellCtrl.ShowMenu()
			}
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			baseURL := cmd.Flag("server").Value.String()
			shellCtrl = NewShellController(baseURL)
		},
	}
)

func init() {
}
