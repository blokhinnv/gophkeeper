package item

import (
	"github.com/spf13/cobra"

	"gophkeeper/internal/client/commands/item/upsert"
	"gophkeeper/internal/client/service"
)

var (
	storageService = service.NewStorageService("http://localhost:8080")
	// ItemCmd represents the item command
	ItemCmd = &cobra.Command{
		Use:   "item",
		Short: "...",
	}
)

func init() {
	ItemCmd.PersistentFlags().String("token", "", "jwt token")
	ItemCmd.PersistentFlags().StringP("type", "t", "", "type of entity to be added")

	ItemCmd.AddCommand(upsert.UpsertCmd)
	ItemCmd.AddCommand(deleteCmd)
}
