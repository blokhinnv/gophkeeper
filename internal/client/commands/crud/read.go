package crud

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/blokhinnv/gophkeeper/internal/server/models"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "read command",
	Long: `The readCmd command retrieves all documents from a specified collection.
It accepts flags to decrypt the data from an encrypted file.
The result is returned as a JSON string.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		key := cmd.Flag("key").Value.String()
		file := cmd.Flag("file").Value.String()
		collectionName, err := models.NewCollectionName(cmd.Flag("collection").Value.String())
		if err != nil {
			fmt.Println(err)
			return err
		}

		decrypted, err := encryptService.FromEncryptedFile(file, key)
		if err != nil {
			fmt.Println(err)
			return err
		}

		res := storageService.GetAll(collectionName, decrypted)
		resJSON, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("Result: %s\n", resJSON)
		return nil
	},
}

func init() {
	readCmd.PersistentFlags().StringP("file", "f", "", "filename to load synced data from")
	readCmd.PersistentFlags().StringP("key", "k", "", "key for data decryption")

	for _, flag := range []string{"file", "key"} {
		readCmd.MarkPersistentFlagRequired(flag)
	}
}
