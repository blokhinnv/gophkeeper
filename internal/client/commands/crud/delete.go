package crud

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/blokhinnv/gophkeeper/internal/server/models"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete command",
	Long: `The 'delete' command deletes a single record from the specified collection in the remote storage service.
It requires a valid authentication token and the record ID to be deleted as a flag.
Provide the name of the collection to delete the record from as an argument.
The command will construct a request body using the provided record ID, and send the DELETE request to the remote service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cmd.Flag("token").Value.String()
		id := cmd.Flag("id").Value.String()
		collectionName, err := models.NewCollectionName(cmd.Flag("collection").Value.String())
		if err != nil {
			fmt.Println(err)
			return err
		}
		body := fmt.Sprintf(`{"record_id": "%v"}`, id)
		msg, err := storageService.Delete(body, collectionName, token)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(msg)
		return nil
	},
}

func init() {
	deleteCmd.PersistentFlags().String("id", "", "id of a record to delete")
	deleteCmd.PersistentFlags().String("token", "t", "user's jwt token")
	for _, flag := range []string{"id", "token"} {
		deleteCmd.MarkPersistentFlagRequired(flag)
	}
}
