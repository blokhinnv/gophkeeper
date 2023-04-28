package crud

import (
	"github.com/spf13/cobra"

	"github.com/blokhinnv/gophkeeper/internal/client/commands/crud/upsert"
	"github.com/blokhinnv/gophkeeper/internal/client/service"
)

var (
	storageService = service.NewStorageService("http://localhost:8080")
	encryptService = service.NewEncryptService()
	// CRUDCmd represents the CRUD command
	CRUDCmd = &cobra.Command{
		Use:   "crud",
		Short: "a command for crud operations",
		Long:  `A parent command for a add, delete and upsert.`,
	}
)

func init() {
	CRUDCmd.PersistentFlags().StringP("collection", "c", "", "a collection to work with")
	CRUDCmd.MarkPersistentFlagRequired("collection")

	CRUDCmd.AddCommand(readCmd)
	CRUDCmd.AddCommand(upsert.UpsertCmd)
	CRUDCmd.AddCommand(deleteCmd)
}
