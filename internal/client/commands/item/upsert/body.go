package upsert

import (
	"encoding/json"
	"fmt"
	"gophkeeper/internal/server/models"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func metadataFromFlags(flagsMetadata MetadataSlice) (models.Metadata, error) {
	md := make(models.Metadata)
	for _, v := range flagsMetadata {
		kv := strings.Split(v, ";")
		if len(kv) != 2 {
			return nil, fmt.Errorf("wrong metadata %v", kv)
		}
		md[kv[0]] = kv[1]
	}
	return md, nil
}

func getBody(flags *UpsertFlags, recordType, recordIDHex string) (string, error) {
	md, err := metadataFromFlags(flags.Metadata)
	if err != nil {
		return "", err
	}

	recordID, err := primitive.ObjectIDFromHex(recordIDHex)
	if err != nil {
		return "", err
	}

	var body any

	switch string(recordType) {
	case models.TextCollection:
		body = &models.TextRecord{
			Data:     flags.TextInfo,
			Metadata: md,
			RecordID: recordID,
		}
	case models.BinaryCollection:
		body = &models.BinaryRecord{Data: flags.BinaryInfo, Metadata: md, RecordID: recordID}
	case models.CardCollection:
		body = &models.CardRecord{Data: flags.CardInfo, Metadata: md, RecordID: recordID}
	case models.CredentialsCollection:
		body = &models.CredentialRecord{
			Data:     flags.CredentialInfo,
			Metadata: md,
			RecordID: recordID,
		}
	default:
		return "", fmt.Errorf("unknown record type")
	}

	bodyEncoded, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	return string(bodyEncoded), nil
}
