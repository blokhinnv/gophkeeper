package auth

import (
	"fmt"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login",
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
		tok, err := authService.Auth(username, password)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Token: ", tok)
	},
}

func init() {
	AuthCmd.AddCommand(loginCmd)
}
