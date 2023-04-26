package upsert

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addCmd represents the login command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add",
	Long:  `...`,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := cmd.Flags().GetString("token")
		if err != nil {
			fmt.Println(err)
			return
		}
		recordType, err := cmd.Flags().GetString("type")
		if err != nil {
			fmt.Println(err)
			return
		}

		body, err := getBody(&cmdFlags, recordType, "000000000000000000000000")
		if err != nil {
			fmt.Println(err)
			return
		}
		err = storageService.Add(body, recordType, token)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	UpsertCmd.AddCommand(addCmd)
}
