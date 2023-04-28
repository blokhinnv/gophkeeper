package upsert

import (
	"fmt"

	"github.com/spf13/cobra"

	"gophkeeper/internal/server/models"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update command",
	Long: `The update command updates an existing record in a specified collection.
It requires a valid token, a collection name, and the id of the record to be updated.
It also expects a valid JSON body that contains the updated information for the record.`,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := cmd.Flags().GetString("token")
		if err != nil {
			fmt.Println(err)
			return
		}
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			fmt.Println(err)
			return
		}
		collectionName, err := models.NewCollection(cmd.Flag("collection").Value.String())
		if err != nil {
			fmt.Println(err)
			return
		}

		body, err := getBody(&cmdFlags, collectionName, id)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = storageService.Update(body, collectionName, token)
		if err != nil {
			fmt.Println(err)
			return
		}

	},
}

func init() {
	UpsertCmd.AddCommand(updateCmd)
	updateCmd.PersistentFlags().String("id", "", "id of a record to update")
	updateCmd.MarkPersistentFlagRequired("id")

}
