package upsert

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/blokhinnv/gophkeeper/internal/server/models"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update command",
	Long: `The update command updates an existing record in a specified collection.
It requires a valid token, a collection name, and the id of the record to be updated.
It also expects a valid JSON body that contains the updated information for the record.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cmd.Flag("token").Value.String()
		id := cmd.Flag("id").Value.String()
		collectionName, err := models.NewCollectionName(cmd.Flag("collection").Value.String())
		if err != nil {
			fmt.Println(err)
			return err
		}

		body, err := getBody(&cmdFlags, collectionName, id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		msg, err := storageService.Update(body, collectionName, token)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(msg)
		return nil

	},
}

func init() {
	UpsertCmd.AddCommand(updateCmd)
	updateCmd.PersistentFlags().String("id", "", "id of a record to update")
	updateCmd.MarkPersistentFlagRequired("id")

}
