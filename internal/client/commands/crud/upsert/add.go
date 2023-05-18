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
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cmd.Flag("token").Value.String()
		collectionName, err := models.NewCollectionName(cmd.Flag("collection").Value.String())
		if err != nil {
			fmt.Println(err)
			return err
		}

		body, err := getBody(&cmdFlags, collectionName, "000000000000000000000000")
		if err != nil {
			fmt.Println(err)
			return err
		}
		msg, err := storageService.Add(body, collectionName, token)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(msg)
		return nil
	},
}

func init() {
	UpsertCmd.AddCommand(addCmd)
}
