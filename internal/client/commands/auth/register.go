package auth

import (
	"fmt"

	"github.com/spf13/cobra"
)

// registerCmd represents the login command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "register",
	Long:  `...`,
	Run: func(cmd *cobra.Command, args []string) {
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			fmt.Println(err)
			return
		}
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			fmt.Println(err)
			return
		}
		err = authService.Register(username, password)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	AuthCmd.AddCommand(registerCmd)
}
