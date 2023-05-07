package service

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	clientModels "github.com/blokhinnv/gophkeeper/internal/client/models"
	srvrModels "github.com/blokhinnv/gophkeeper/internal/server/models"
)

func TestEncryptService_ToEncryptedFile(t *testing.T) {
	// create a temporary file
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	data := &clientModels.SyncResponse{
		Text: []srvrModels.TextRecord{
			{
				RecordID: primitive.NewObjectID(),
				Username: "blokhinnv",
				Data:     srvrModels.TextInfo("some data..."),
			},
		},
	}
	// create an instance of the encryptService
	s := NewEncryptService()

	// call ToEncryptedFile with a password
	password := "password"
	err = s.ToEncryptedFile(data, tmpfile.Name(), password)

	// check that no errors occurred
	assert.NoError(t, err)

	// read the file contents
	fileContents, err := os.ReadFile(tmpfile.Name())
	assert.NoError(t, err)

	// check that the file contents are not equal to the original data
	c, err := json.Marshal(fileContents)
	assert.NotEqual(t, string(fileContents), string(c))
	assert.NoError(t, err)
}

func TestEncryptService_FromEncryptedFile(t *testing.T) {
	// create a temporary file
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	data := &clientModels.SyncResponse{
		Text: []srvrModels.TextRecord{
			{
				RecordID: primitive.NewObjectID(),
				Username: "blokhinnv",
				Data:     srvrModels.TextInfo("some data..."),
			},
		},
	}
	// create an instance of the encryptService
	s := NewEncryptService()

	// call ToEncryptedFile with a password
	password := "password"
	err = s.ToEncryptedFile(data, tmpfile.Name(), password)
	assert.NoError(t, err)

	resp, err := s.FromEncryptedFile(tmpfile.Name(), password)
	// Assert that the response is not nil and there is no error
	assert.NotNil(t, resp)
	assert.NoError(t, err)
	// Assert that the response contains the decrypted data
	assert.Equal(t, "some data...", resp.Text[0].Data)
}
