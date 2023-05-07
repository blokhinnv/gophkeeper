package sync

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/blokhinnv/gophkeeper/internal/client/commands/cotesting"
	"github.com/blokhinnv/gophkeeper/internal/client/service/mock"
	"github.com/blokhinnv/gophkeeper/internal/server/models"
)

func init() {
	SyncCmd.PersistentFlags().StringP("server", "s", "https://localhost:8080", "server addr")
}

func TestSyncCommand(t *testing.T) {
	SyncCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		syncService = mock.NewMockSyncService(mockCtrl)
		encryptService = mock.NewMockEncryptService(mockCtrl)
		syncService.(*mock.MockSyncService).EXPECT().
			Sync(gomock.Eq("sometoken"), gomock.Eq([]models.CollectionName{"text"})).
			AnyTimes().
			Return(nil, nil)

		syncService.(*mock.MockSyncService).EXPECT().
			Sync(gomock.Eq("sometoken"), gomock.Eq([]models.CollectionName{"binary"})).
			AnyTimes().
			Return(nil, fmt.Errorf("unable to sync"))

		encryptService.(*mock.MockEncryptService).EXPECT().
			ToEncryptedFile(gomock.Any(), gomock.Eq("fname"), gomock.Eq("somekey")).
			AnyTimes().
			Return(nil)
	}

	rootCmd := SyncCmd
	t.Run("ok", func(t *testing.T) {
		defer rootCmd.ResetFlags()
		err := cotesting.ExecuteCommandC(
			rootCmd,
			"sync",
			"--token=sometoken",
			"--file=fname",
			"--key=somekey",
			"--collection=text",
		)
		assert.NoError(t, err)
	})
	t.Run("no_sync", func(t *testing.T) {
		defer rootCmd.ResetFlags()
		err := cotesting.ExecuteCommandC(
			rootCmd,
			"sync",
			"--token=sometoken",
			"--file=fname",
			"--key=somekey",
			"--collection=binary",
		)
		assert.Error(t, err)
	})
	t.Run("bad_collection", func(t *testing.T) {
		defer rootCmd.ResetFlags()
		err := cotesting.ExecuteCommandC(
			rootCmd,
			"sync",
			"--token=sometoken",
			"--file=fname",
			"--key=somekey",
			"--collection=somecoll1",
		)
		assert.Error(t, err)
	})
}
