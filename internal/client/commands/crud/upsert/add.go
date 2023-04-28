package upsert

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/blokhinnv/gophkeeper/internal/server/models"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add command",
	Long: `The add command allows users to add new records to a specified collection.
It requires a valid JWT token for authorization and accepts various flags for
different types of data. The command constructs the record body and sends it
to the server for storage.`,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := cmd.Flags().GetString("token")
		if err != nil {
			fmt.Println(err)
			return
		}
		collectionName, err := models.NewCollection(cmd.Flag("collection").Value.String())
		if err != nil {
			fmt.Println(err)
			return
		}

		body, err := getBody(&cmdFlags, collectionName, "000000000000000000000000")
		if err != nil {
			fmt.Println(err)
			return
		}
		err = storageService.Add(body, collectionName, token)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	UpsertCmd.AddCommand(addCmd)
}
