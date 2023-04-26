package upsert

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addCmd represents the login command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update",
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

		body, err := getBody(&cmdFlags, recordType, id)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = storageService.Update(body, recordType, token)
		if err != nil {
			fmt.Println(err)
			return
		}

	},
}

func init() {
	UpsertCmd.AddCommand(updateCmd)
	updateCmd.PersistentFlags().String("id", "", "id of a record to update")
}
