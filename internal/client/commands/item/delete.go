package item

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addCmd represents the login command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete",
	Long:  `...`,
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
		recordType, err := cmd.Flags().GetString("type")
		if err != nil {
			fmt.Println(err)
			return
		}
		body := fmt.Sprintf(`{"record_id": "%v"}`, id)
		err = storageService.Delete(body, recordType, token)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	deleteCmd.PersistentFlags().String("id", "", "id of a record to delete")
}
