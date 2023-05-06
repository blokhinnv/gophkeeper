package upsert

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/blokhinnv/gophkeeper/internal/server/errors"
	"github.com/blokhinnv/gophkeeper/internal/server/models"
)

// metadataFromFlags converts metadata from flags to models.Metadata
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

// fileToBase64 read file's bytes and encodes them in base64.
func fileToBase64(filename string) (string, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// getBody returns the body for the upsert operation
func getBody(
	flags *UpsertFlags,
	collectionName models.CollectionName,
	recordIDHex string,
) (string, error) {
	md, err := metadataFromFlags(flags.Metadata)
	if err != nil {
		return "", err
	}

	recordID, err := primitive.ObjectIDFromHex(recordIDHex)
	if err != nil {
		return "", err
	}

	var body any

	switch collectionName {
	case models.TextCollection:
		body = &models.TextRecord{
			Data:     flags.TextInfo,
			Metadata: md,
			RecordID: recordID,
		}
	case models.BinaryCollection:
		fname := flags.BinaryInfo.FileName
		content, err := fileToBase64(fname)
		if err != nil {
			return "", err
		}
		data := models.BinaryInfo{
			FileName: fname,
			Content:  content,
		}
		body = &models.BinaryRecord{Data: data, Metadata: md, RecordID: recordID}
	case models.CardCollection:
		body = &models.CardRecord{Data: flags.CardInfo, Metadata: md, RecordID: recordID}
	case models.CredentialsCollection:
		body = &models.CredentialRecord{
			Data:     flags.CredentialInfo,
			Metadata: md,
			RecordID: recordID,
		}
	default:
		return "", fmt.Errorf("%w: %v", errors.ErrUnknownCollection, collectionName)
	}

	bodyEncoded, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	return string(bodyEncoded), nil
}
