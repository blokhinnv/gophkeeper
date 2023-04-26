package upsert

import (
	"gophkeeper/internal/client/service"
	"gophkeeper/internal/server/models"

	"github.com/spf13/cobra"
)

type MetadataSlice = []string

type UpsertFlags struct {
	models.TextInfo
	models.BinaryInfo
	models.CredentialInfo
	models.CardInfo
	Metadata MetadataSlice
}

var (
	cmdFlags       = UpsertFlags{}
	storageService = service.NewStorageService("http://localhost:8080")
	// UpsertCmd represents the item command
	UpsertCmd = &cobra.Command{
		Use:   "upsert",
		Short: "...",
	}
)

func init() {
	UpsertCmd.PersistentFlags().
		StringVar(&cmdFlags.TextInfo, "text", "", "data for a text record")

	UpsertCmd.PersistentFlags().
		StringVar(&cmdFlags.BinaryInfo, "binary-data", "", "data for a binary record")

	UpsertCmd.PersistentFlags().
		StringVar(&cmdFlags.CredentialInfo.Login, "login", "", "data for a credentials record")
	UpsertCmd.PersistentFlags().
		StringVar(&cmdFlags.CredentialInfo.Password, "password", "", "data for a credentials record")

	UpsertCmd.PersistentFlags().
		StringVar(&cmdFlags.CardInfo.CardNumber, "card-number", "", "data for a credentials record")
	UpsertCmd.PersistentFlags().
		StringVar(&cmdFlags.CardInfo.CVV, "cvv", "", "data for a credentials record")
	UpsertCmd.PersistentFlags().
		StringVar(&cmdFlags.CardInfo.ExpirationDate, "expiration-date", "", "data for a credentials record")

	UpsertCmd.PersistentFlags().
		StringSliceVarP(&cmdFlags.Metadata, "meta", "m", []string{}, "semicolor separated metadata values")

}
