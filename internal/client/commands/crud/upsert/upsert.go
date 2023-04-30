// package upsert provides functions to handle upsert operations
package upsert

import (
	"github.com/spf13/cobra"

	"github.com/blokhinnv/gophkeeper/internal/client/service"
	"github.com/blokhinnv/gophkeeper/internal/server/models"
)

// MetadataSlice is a slice of strings to store metadata
type MetadataSlice = []string

// UpsertFlags holds the different types of records and their respective fields,
// along with metadata for the record.
type UpsertFlags struct {
	models.TextInfo
	models.BinaryInfo
	models.CredentialInfo
	models.CardInfo
	Metadata MetadataSlice
}

var (
	cmdFlags       = UpsertFlags{}
	storageService service.StorageService
	// UpsertCmd represents the item command
	UpsertCmd = &cobra.Command{
		Use:   "upsert",
		Short: "upsert command",
		Long:  "A parent command for add and update.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			baseURL := cmd.Flag("server").Value.String()
			storageService = service.NewStorageService(baseURL)
		},
	}
)

func init() {
	UpsertCmd.PersistentFlags().String("token", "", "user's jwt token")
	UpsertCmd.MarkPersistentFlagRequired("token")

	UpsertCmd.PersistentFlags().
		StringVar(&cmdFlags.TextInfo, "text", "", "data for a text record")

	UpsertCmd.PersistentFlags().
		StringVar(&cmdFlags.BinaryInfo.FileName, "file", "", "path to file which will be stored")

	UpsertCmd.PersistentFlags().
		StringVar(&cmdFlags.CredentialInfo.Login, "login", "", "data for a credentials record")
	UpsertCmd.PersistentFlags().
		StringVar(&cmdFlags.CredentialInfo.Password, "password", "", "data for a credentials record")
	UpsertCmd.MarkFlagsRequiredTogether("login", "password")

	UpsertCmd.PersistentFlags().
		StringVar(&cmdFlags.CardInfo.CardNumber, "card-number", "", "data for a credentials record")
	UpsertCmd.PersistentFlags().
		StringVar(&cmdFlags.CardInfo.CVV, "cvv", "", "data for a credentials record")
	UpsertCmd.PersistentFlags().
		StringVar(&cmdFlags.CardInfo.ExpirationDate, "expiration-date", "", "data for a credentials record")
	UpsertCmd.MarkFlagsRequiredTogether("card-number", "cvv", "expiration-date")

	UpsertCmd.PersistentFlags().
		StringSliceVarP(&cmdFlags.Metadata, "meta", "m", []string{}, "semicolor separated metadata values")

}
