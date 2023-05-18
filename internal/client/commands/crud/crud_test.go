package crud

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/blokhinnv/gophkeeper/internal/client/commands/cotesting"
	clientModels "github.com/blokhinnv/gophkeeper/internal/client/models"
	"github.com/blokhinnv/gophkeeper/internal/client/service/mock"
	srvrModels "github.com/blokhinnv/gophkeeper/internal/server/models"
)

func init() {
	CRUDCmd.PersistentFlags().StringP("server", "s", "https://localhost:8080", "server addr")
}

func TestReadCommand(t *testing.T) {
	CRUDCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		storageService = mock.NewMockStorageService(mockCtrl)
		encryptService = mock.NewMockEncryptService(mockCtrl)

		r := &clientModels.SyncResponse{
			Text: []srvrModels.TextRecord{
				{Data: "some text"},
			},
		}

		encryptService.(*mock.MockEncryptService).EXPECT().
			FromEncryptedFile(gomock.Eq("badfname"), gomock.Eq("badkey")).
			AnyTimes().
			Return(r, fmt.Errorf("bad file"))

		encryptService.(*mock.MockEncryptService).EXPECT().
			FromEncryptedFile(gomock.Eq("fname"), gomock.Eq("correctkey")).
			AnyTimes().
			Return(r, nil)

		storageService.(*mock.MockStorageService).EXPECT().
			GetAll(srvrModels.CollectionName("text"), r).
			AnyTimes().
			Return(r.Text)
	}

	rootCmd := CRUDCmd
	t.Run("bad_collection", func(t *testing.T) {
		err := cotesting.ExecuteCommandC(
			rootCmd,
			"read",
			"--key=somekey",
			"--file=fname",
			"--collection=badcollection",
		)
		assert.Error(t, err)
	})
	t.Run("bad_file", func(t *testing.T) {
		err := cotesting.ExecuteCommandC(
			rootCmd,
			"read",
			"--key=badkey",
			"--file=badfname",
			"--collection=badcollection",
		)
		assert.Error(t, err)
	})
	t.Run("ok", func(t *testing.T) {
		err := cotesting.ExecuteCommandC(
			rootCmd,
			"read",
			"--key=correctkey",
			"--file=fname",
			"--collection=text",
		)
		assert.NoError(t, err)
	})
}
func TestDeleteCommand(t *testing.T) {
	CRUDCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		storageService = mock.NewMockStorageService(mockCtrl)
		encryptService = mock.NewMockEncryptService(mockCtrl)

		storageService.(*mock.MockStorageService).EXPECT().
			Delete(`{"record_id": "1234"}`, srvrModels.CollectionName("text"), "sometoken").
			AnyTimes().
			Return("ok", nil)

		storageService.(*mock.MockStorageService).EXPECT().
			Delete(`{"record_id": "0000"}`, srvrModels.CollectionName("text"), "sometoken").
			AnyTimes().
			Return("", fmt.Errorf("not found"))
	}

	rootCmd := CRUDCmd
	t.Run("bad_collection", func(t *testing.T) {
		err := cotesting.ExecuteCommandC(
			rootCmd,
			"delete",
			"--token=sometoken",
			"--id=111111",
			"--collection=badcollection",
		)
		assert.Error(t, err)
	})
	t.Run("not_found", func(t *testing.T) {
		err := cotesting.ExecuteCommandC(
			rootCmd,
			"delete",
			"--token=sometoken",
			"--id=0000",
			"--collection=text",
		)
		assert.Error(t, err)
	})
	t.Run("ok", func(t *testing.T) {
		err := cotesting.ExecuteCommandC(
			rootCmd,
			"delete",
			"--token=sometoken",
			"--id=1234",
			"--collection=text",
		)
		assert.NoError(t, err)
	})
}
