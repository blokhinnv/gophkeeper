package shell

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ShellCmd represents the auth command
var ShellCmd = &cobra.Command{
	Use:   "shell",
	Short: "...",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("shell called")
	},
}

func init() {
}
