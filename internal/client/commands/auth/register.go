package auth

import (
	"fmt"

	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "register command",
	Long: `The register command is used to register a new user with the authentication service.
It requires a username and password flag to be set.
Upon successful registration, the command will return no output.
If the registration fails, an error message will be printed to the console.`,
	Run: func(cmd *cobra.Command, args []string) {
		username := cmd.Flag("username").Value.String()
		password := cmd.Flag("password").Value.String()
		err := authService.Register(username, password)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	AuthCmd.AddCommand(registerCmd)
}
