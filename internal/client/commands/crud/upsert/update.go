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
	Run: func(cmd *cobra.Command, args []string) {
		token := cmd.Flag("token").Value.String()
		id := cmd.Flag("id").Value.String()
		collectionName, err := models.NewCollectionName(cmd.Flag("collection").Value.String())
		if err != nil {
			fmt.Println(err)
			return
		}

		body, err := getBody(&cmdFlags, collectionName, id)
		if err != nil {
			fmt.Println(err)
			return
		}
		msg, err := storageService.Update(body, collectionName, token)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(msg)

	},
}

func init() {
	UpsertCmd.AddCommand(updateCmd)
	updateCmd.PersistentFlags().String("id", "", "id of a record to update")
	updateCmd.MarkPersistentFlagRequired("id")

}
