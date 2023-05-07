package models

import "github.com/blokhinnv/gophkeeper/internal/server/models"

// SyncResponse defines the response from the SyncService Sync method.
type SyncResponse struct {
	Text       []models.TextRecord
	Binary     []models.BinaryRecord
	Card       []models.CardRecord
	Credential []models.CredentialRecord
}
