// Package crud provides implementations of CRUD CLI-commands.
package crud

import (
	"github.com/spf13/cobra"

	"github.com/blokhinnv/gophkeeper/internal/client/commands/crud/upsert"
	"github.com/blokhinnv/gophkeeper/internal/client/service"
)

var (
	// storageService is a storage service used for a command implementation.
	storageService service.StorageService
	// storageService is a encryption service used for a command implementation.
	encryptService service.EncryptService
	// CRUDCmd represents the CRUD command.
	CRUDCmd = &cobra.Command{
		Use:   "crud",
		Short: "a command for crud operations",
		Long:  `A parent command for a add, delete and upsert.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			baseURL := cmd.Flag("server").Value.String()
			storageService = service.NewStorageService(baseURL)
			encryptService = service.NewEncryptService()
		},
	}
)

func init() {
	CRUDCmd.PersistentFlags().StringP("collection", "c", "", "a collection to work with")
	CRUDCmd.MarkPersistentFlagRequired("collection")
	CRUDCmd.AddCommand(readCmd, deleteCmd, upsert.UpsertCmd)
}
