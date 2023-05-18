package auth

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/blokhinnv/gophkeeper/internal/client/commands/cotesting"
	"github.com/blokhinnv/gophkeeper/internal/client/service/mock"
)

func init() {
	AuthCmd.PersistentFlags().StringP("server", "s", "https://localhost:8080", "server addr")
}

func TestLoginCommand(t *testing.T) {
	AuthCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		authService = mock.NewMockAuthService(mockCtrl)
		authService.(*mock.MockAuthService).EXPECT().
			Auth(gomock.Eq("someuser"), gomock.Eq("correctpwd")).
			AnyTimes().
			Return("some token", nil)
		authService.(*mock.MockAuthService).EXPECT().
			Auth(gomock.Eq("someuser"), gomock.Eq("wrongpwd")).
			AnyTimes().
			Return("", fmt.Errorf("Bad credentials"))
	}

	rootCmd := AuthCmd
	t.Run("ok", func(t *testing.T) {
		err := cotesting.ExecuteCommandC(
			rootCmd,
			"login",
			"--username=someuser",
			"--password=correctpwd",
		)
		assert.NoError(t, err)
	})
	t.Run("bad", func(t *testing.T) {
		err := cotesting.ExecuteCommandC(
			rootCmd,
			"login",
			"--username=someuser",
			"--password=wrongpwd",
		)
		assert.Error(t, err)
	})
}

func TestRegisterCommand(t *testing.T) {
	AuthCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		authService = mock.NewMockAuthService(mockCtrl)
		authService.(*mock.MockAuthService).EXPECT().
			Register(gomock.Eq("someuser"), gomock.Eq("correctpwd")).
			AnyTimes().
			Return(nil)
		authService.(*mock.MockAuthService).EXPECT().
			Register(gomock.Eq("someuserexisted"), gomock.Eq("correctpwd")).
			AnyTimes().
			Return(fmt.Errorf("Bad credentials"))
	}
	rootCmd := AuthCmd
	t.Run("ok", func(t *testing.T) {
		err := cotesting.ExecuteCommandC(
			rootCmd,
			"register",
			"--username=someuser",
			"--password=correctpwd",
		)
		assert.NoError(t, err)
	})
	t.Run("bad", func(t *testing.T) {
		err := cotesting.ExecuteCommandC(
			rootCmd,
			"register",
			"--username=someuserexisted",
			"--password=correctpwd",
		)
		assert.Error(t, err)
	})
}
