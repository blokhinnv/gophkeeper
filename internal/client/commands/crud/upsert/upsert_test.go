package upsert

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/blokhinnv/gophkeeper/internal/client/commands/cotesting"
	"github.com/blokhinnv/gophkeeper/internal/client/service/mock"
	"github.com/blokhinnv/gophkeeper/internal/server/models"
)

func init() {
	UpsertCmd.PersistentFlags().StringP("server", "s", "https://localhost:8080", "server addr")
	UpsertCmd.PersistentFlags().StringP("collection", "c", "", "a collection to work with")
	UpsertCmd.PersistentFlags().String("id", "", "id of a record to update")
}

func TestAddCommand(t *testing.T) {
	UpsertCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		storageService = mock.NewMockStorageService(mockCtrl)
		storageService.(*mock.MockStorageService).EXPECT().
			Add(gomock.Any(), gomock.Eq(models.CollectionName("text")), gomock.Eq("sometoken")).
			AnyTimes().
			Return("ok", nil)
	}

	rootCmd := UpsertCmd
	t.Run("ok", func(t *testing.T) {
		err := cotesting.ExecuteCommandC(
			rootCmd,
			"add",
			"--token=sometoken",
			"--collection=text",
			"--text=sometext",
		)
		assert.NoError(t, err)
	})
	t.Run("bad_collection", func(t *testing.T) {
		err := cotesting.ExecuteCommandC(
			rootCmd,
			"add",
			"--token=sometoken",
			"--collection=qwerty",
			"--text=sometext",
		)
		assert.Error(t, err)
	})
}

func TestUpdateCommand(t *testing.T) {
	UpsertCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		storageService = mock.NewMockStorageService(mockCtrl)
		storageService.(*mock.MockStorageService).EXPECT().
			Update(gomock.Any(), gomock.Eq(models.CollectionName("text")), gomock.Eq("sometoken")).
			AnyTimes().
			Return("ok", nil)
	}

	rootCmd := UpsertCmd
	t.Run("ok", func(t *testing.T) {
		err := cotesting.ExecuteCommandC(
			rootCmd,
			"update",
			"--token=sometoken",
			"--collection=text",
			"--text=sometext",
			"--id=6457e99ec51d35bd689f2f5b",
		)
		assert.NoError(t, err)
	})
	t.Run("bad_collection", func(t *testing.T) {
		err := cotesting.ExecuteCommandC(
			rootCmd,
			"update",
			"--token=sometoken",
			"--collection=qwerty",
			"--text=sometext",
			"--id=6457e99ec51d35bd689f2f5b",
		)
		assert.Error(t, err)
	})
	t.Run("bad_body", func(t *testing.T) {
		err := cotesting.ExecuteCommandC(
			rootCmd,
			"update",
			"--token=sometoken",
			"--collection=qwerty",
			"--id=6457e99ec51d35bd689f2f5b",
		)
		assert.Error(t, err)
	})
	t.Run("bad_id", func(t *testing.T) {
		err := cotesting.ExecuteCommandC(
			rootCmd,
			"update",
			"--token=sometoken",
			"--collection=qwerty",
			"--id=1111",
		)
		assert.Error(t, err)
	})
}
