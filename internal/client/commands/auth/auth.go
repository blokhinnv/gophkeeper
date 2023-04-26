package auth

import (
	"gophkeeper/internal/client/service"

	"github.com/spf13/cobra"
)

// AuthCmd represents the auth command
var (
	authService = service.NewAuthService("http://localhost:8080")
	AuthCmd     = &cobra.Command{
		Use:   "auth",
		Short: "...",
	}
)

func init() {
	AuthCmd.PersistentFlags().StringP("username", "u", "", "username to authorize")
	AuthCmd.PersistentFlags().StringP("password", "p", "", "password to authorize")
}
