package auth

import (
	"github.com/spf13/cobra"

	"github.com/blokhinnv/gophkeeper/internal/client/service"
)

// AuthCmd represents the auth command
var (
	authService service.AuthService
	AuthCmd     = &cobra.Command{
		Use:   "auth",
		Short: "authorization and registration commands",
		Long:  "A parent command for login and register.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			baseURL := cmd.Flag("server").Value.String()
			authService = service.NewAuthService(baseURL)
		},
	}
)

func init() {
	AuthCmd.PersistentFlags().StringP("username", "u", "", "username to authorize")
	AuthCmd.PersistentFlags().StringP("password", "p", "", "password to authorize")

	for _, flag := range []string{"username", "password"} {
		AuthCmd.MarkPersistentFlagRequired(flag)
	}
}
