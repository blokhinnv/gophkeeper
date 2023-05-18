package shell

import (
	"encoding/json"

	"github.com/blokhinnv/gophkeeper/internal/server/models"
)

// getBody function is responsible for generating the body of a request to be sent to the server.
// It takes in a models.Collection parameter and a boolean value indicating if request ID is required.
// It returns a string containing the encoded body and an error if one occurs.
func getBody(collection models.CollectionName, requestID bool) (string, error) {
	var data any

	var recordIDHex string = "000000000000000000000000"
	if requestID {
		recordIDHex = promptText("Record id: ")
	}
	recordID, err := models.ObjectIDFromString(recordIDHex)
	if err != nil {
		return "", err
	}

	switch collection {
	case models.TextCollection:
		data = promptText("Text data")
	case models.BinaryCollection:
		data = promptText("Binary data")
	case models.CredentialsCollection:
		login := promptText("Login")
		pwd := promptText("Password")
		data = models.CredentialInfo{Login: login, Password: pwd}
	case models.CardCollection:
		num := promptText("Card number")
		cvv := promptText("CVV")
		exp := promptText("Expiration date")
		data = models.CardInfo{CardNumber: num, CVV: cvv, ExpirationDate: exp}
	}
	md := getMetadata()
	body := &models.UntypedRecord{
		UntypedRecordContent: models.UntypedRecordContent{
			Data:     data,
			Metadata: md,
		},
		RecordID: recordID,
	}
	bodyEncoded, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	return string(bodyEncoded), nil
}
