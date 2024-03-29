package auth

import (
	"fmt"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login",
	Long: `The loginCmd command represents the login functionality, used for user authorization.
The command takes a username and password as arguments and returns a token,
which can be used for subsequent authenticated requests.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		username := cmd.Flag("username").Value.String()
		password := cmd.Flag("password").Value.String()
		tok, err := authService.Auth(username, password)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("Token: ", tok)
		return nil
	},
}

func init() {
	AuthCmd.AddCommand(loginCmd)
}
