package auth

import (
	"github.com/spf13/cobra"

	"gophkeeper/internal/client/service"
)

// AuthCmd represents the auth command
var (
	authService = service.NewAuthService("http://localhost:8080")
	AuthCmd     = &cobra.Command{
		Use:   "auth",
		Short: "authorization and registration commands",
		Long:  "A parent command for login and register.",
	}
)

func init() {
	AuthCmd.PersistentFlags().StringP("username", "u", "", "username to authorize")
	AuthCmd.PersistentFlags().StringP("password", "p", "", "password to authorize")

	for _, flag := range []string{"username", "password"} {
		AuthCmd.MarkPersistentFlagRequired(flag)
	}
}
