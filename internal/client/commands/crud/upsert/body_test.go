package upsert

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/blokhinnv/gophkeeper/internal/server/errors"
	"github.com/blokhinnv/gophkeeper/internal/server/models"
)

func TestMetadataFromFlags(t *testing.T) {
	tests := []struct {
		name        string
		flags       MetadataSlice
		expectedMd  models.Metadata
		expectedErr error
	}{
		{
			name: "valid metadata",
			flags: MetadataSlice{
				"key1;value1",
				"key2;value2",
			},
			expectedMd: models.Metadata{
				"key1": "value1",
				"key2": "value2",
			},
			expectedErr: nil,
		},
		{
			name: "invalid metadata",
			flags: MetadataSlice{
				"key1",
				"key2:value2",
			},
			expectedMd:  nil,
			expectedErr: fmt.Errorf("wrong metadata [key1]"),
		},
		{
			name:        "empty metadata",
			flags:       MetadataSlice{},
			expectedMd:  models.Metadata{},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md, err := metadataFromFlags(tt.flags)
			if !reflect.DeepEqual(md, tt.expectedMd) {
				t.Errorf("metadataFromFlags() got md = %v, expected %v", md, tt.expectedMd)
			}
			if err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf(
					"metadataFromFlags() got err = %v, expected %v",
					err.Error(),
					tt.expectedErr,
				)
			}
		})
	}
}

func TestFileToBase64(t *testing.T) {
	tests := []struct {
		name           string
		filename       string
		expectedResult string
		expectedError  error
	}{
		{
			name:           "valid file",
			filename:       "sample.txt",
			expectedResult: "aGVsbG8sIGdv",
			expectedError:  nil,
		},
		{
			name:           "invalid file",
			filename:       "nonexistent.txt",
			expectedResult: "",
			expectedError: fmt.Errorf(
				"open nonexistent.txt: The system cannot find the file specified.",
			),
		},
		{
			name:           "empty file",
			filename:       "empty.txt",
			expectedResult: "",
			expectedError:  nil,
		},
	}

	err := os.WriteFile("sample.txt", []byte("hello, go"), 0644)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	defer os.Remove("sample.txt")

	err = os.WriteFile("empty.txt", []byte(""), 0644)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	defer os.Remove("empty.txt")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fileToBase64(tt.filename)

			if result != tt.expectedResult {
				t.Errorf("fileToBase64() got %v, expected %v", result, tt.expectedResult)
			}

			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("fileToBase64() got %v, expected %v", err, tt.expectedError)
			}
		})
	}
}

func TestGetBody(t *testing.T) {
	err := os.WriteFile("sample.txt", []byte("hello, go"), 0644)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	defer os.Remove("sample.txt")

	tests := []struct {
		name           string
		flags          *UpsertFlags
		collectionName models.CollectionName
		recordIDHex    string
		expectedBody   string
		expectedError  error
	}{
		{
			name: "valid text collection",
			flags: &UpsertFlags{
				TextInfo: "test text",
				Metadata: MetadataSlice{"key1;value1", "key2;value2"},
			},
			collectionName: models.TextCollection,
			recordIDHex:    "1234567890abcdef12345678",
			expectedBody:   `{"Data":"test text","Metadata":{"key1":"value1","key2":"value2"},"record_id":"1234567890abcdef12345678"}`,
			expectedError:  nil,
		},
		{
			name: "valid binary collection",
			flags: &UpsertFlags{
				BinaryInfo: models.BinaryInfo{
					FileName: "sample.txt",
				},
				Metadata: MetadataSlice{"key1;value1", "key2;value2"},
			},
			collectionName: models.BinaryCollection,
			recordIDHex:    "1234567890abcdef12345678",
			expectedBody:   `{"Data":{"FileName":"sample.txt","Content":"aGVsbG8sIGdv"},"Metadata":{"key1":"value1","key2":"value2"},"record_id":"1234567890abcdef12345678"}`,
			expectedError:  nil,
		},
		{
			name: "valid card collection",
			flags: &UpsertFlags{
				CardInfo: models.CardInfo{
					CardNumber:     "1111 1111 1111 1111",
					CVV:            "123",
					ExpirationDate: "12/12",
				},
				Metadata: MetadataSlice{"key1;value1", "key2;value2"},
			},
			collectionName: models.CardCollection,
			recordIDHex:    "1234567890abcdef12345678",
			expectedBody:   `{"Data":{"CardNumber":"1111 1111 1111 1111","CVV":"123", "ExpirationDate":"12/12"},"Metadata":{"key1":"value1","key2":"value2"},"record_id":"1234567890abcdef12345678"}`,
			expectedError:  nil,
		},
		{
			name: "valid credentials collection",
			flags: &UpsertFlags{
				CredentialInfo: models.CredentialInfo{
					Login:    "user1",
					Password: "password1",
				},
				Metadata: MetadataSlice{"key1;value1", "key2;value2"},
			},
			collectionName: models.CredentialsCollection,
			recordIDHex:    "1234567890abcdef12345678",
			expectedBody:   `{"Data":{"Login":"user1","Password":"password1"},"Metadata":{"key1":"value1","key2":"value2"},"record_id":"1234567890abcdef12345678"}`,
			expectedError:  nil,
		},
		{
			name: "unknown collection",
			flags: &UpsertFlags{
				TextInfo: "test text",
				Metadata: MetadataSlice{"key1;value1", "key2;value2"},
			},
			collectionName: "unknown",
			recordIDHex:    "1234567890abcdef12345678",
			expectedBody:   "",
			expectedError:  errors.ErrUnknownCollection,
		},
		{
			name: "invalid metadata",
			flags: &UpsertFlags{
				TextInfo: "test text",
				Metadata: MetadataSlice{"key1"},
			},
			collectionName: "unknown",
			recordIDHex:    "1234567890abcdef12345678",
			expectedBody:   "",
			expectedError:  fmt.Errorf("md error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := getBody(tt.flags, tt.collectionName, tt.recordIDHex)
			if (err == nil) != (tt.expectedError == nil) {
				t.Errorf(
					"metadataFromFlags() got err = %v, expectedError %v; expected equal",
					err,
					tt.expectedError,
				)
			}

			if err != nil {
				return
			}

			var o1 any
			var o2 any
			err = json.Unmarshal([]byte(body), &o1)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			err = json.Unmarshal([]byte(tt.expectedBody), &o2)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(o1, o2) {
				t.Errorf(
					"getBody() got body = %v, expected %v",
					body,
					tt.expectedBody,
				)
			}

		})
	}
}
