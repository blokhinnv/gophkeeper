// Package sync provides implementations of sync CLI-commands.
package sync

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/blokhinnv/gophkeeper/internal/client/service"
	"github.com/blokhinnv/gophkeeper/internal/server/models"
)

var (
	// syncService is a sync service used for a command implementation.
	syncService service.SyncService
	// encryptService is a encrypt service used for a command implementation.
	encryptService service.EncryptService
	// SyncCmd represents the sync command
	SyncCmd = &cobra.Command{
		Use:   "sync",
		Short: "sync command",
		Long: `The SyncCmd command performs a synchronization operation between
the client and the remote storage service. It accepts the "token", "key",
and "file" flags to authenticate and encrypt the data, respectively.
It also requires the "collection" flag to be set to a list of collections to sync.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			token := cmd.Flag("token").Value.String()
			key := cmd.Flag("key").Value.String()
			file := cmd.Flag("file").Value.String()
			collectionStringSlice, err := cmd.Flags().GetStringSlice("collection")
			if err != nil {
				fmt.Println(err)
				return err
			}
			collections := make([]models.CollectionName, 0, len(collectionStringSlice))
			// loop through the []string slice and convert each element to a Collection
			for _, s := range collectionStringSlice {
				c, err := models.NewCollectionName(s)
				if err != nil {
					fmt.Println(err)
					return err
				}
				collections = append(collections, c)
			}
			resp, err := syncService.Sync(token, collections)
			if err != nil {
				fmt.Println(err)
				return err
			}
			err = encryptService.ToEncryptedFile(resp, file, key)
			if err != nil {
				fmt.Println(err)
				return err
			}
			return err
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			baseURL := cmd.Flag("server").Value.String()
			syncService = service.NewSyncService(baseURL)
			encryptService = service.NewEncryptService()
		},
	}
)

func init() {
	SyncCmd.PersistentFlags().StringP("token", "t", "", "jwt token")
	SyncCmd.PersistentFlags().StringP("file", "f", "", "filename to save synced data")
	SyncCmd.PersistentFlags().StringP("key", "k", "", "key for data encryption")
	SyncCmd.PersistentFlags().
		StringSliceP(
			"collection",
			"c",
			[]string{
				string(models.CardCollection),
				string(models.TextCollection),
				string(models.CredentialsCollection),
				string(models.BinaryCollection),
			},
			"collections to sync",
		)
	for _, flag := range []string{"token", "file", "key"} {
		SyncCmd.MarkPersistentFlagRequired(flag)
	}
}
