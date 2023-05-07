package service

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	clientModels "github.com/blokhinnv/gophkeeper/internal/client/models"
	srvErrors "github.com/blokhinnv/gophkeeper/internal/server/errors"
	srvrModels "github.com/blokhinnv/gophkeeper/internal/server/models"
)

func TestSyncService_Sync(t *testing.T) {
	baseURL := "https://example.com"
	s := NewSyncService(baseURL)
	client := s.GetClient()
	assert.Equal(t, baseURL, client.HostURL)

	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	t.Run("success", func(t *testing.T) {
		httpmock.Reset()
		collectionNames := []srvrModels.CollectionName{
			srvrModels.TextCollection,
			srvrModels.BinaryCollection,
			srvrModels.CardCollection,
			srvrModels.CredentialsCollection,
		}
		expectedResult := &clientModels.SyncResponse{
			Text: []srvrModels.TextRecord{
				{RecordID: primitive.NewObjectID(), Data: srvrModels.TextInfo("some-text...")},
			},
			Binary: []srvrModels.BinaryRecord{
				{
					RecordID: primitive.NewObjectID(),
					Data:     srvrModels.BinaryInfo{FileName: "test.test", Content: "cXdlcXdld3Fl"},
				},
			},
			Card: []srvrModels.CardRecord{
				{
					RecordID: primitive.NewObjectID(),
					Data: srvrModels.CardInfo{
						CardNumber:     "1234 1234 1234 1234",
						CVV:            "234",
						ExpirationDate: "01/12",
					},
				},
			},
			Credential: []srvrModels.CredentialRecord{
				{
					RecordID: primitive.NewObjectID(),
					Data: srvrModels.CredentialInfo{
						Login:    "some-login",
						Password: "some-password",
					},
				},
			},
		}
		responderText, err := httpmock.NewJsonResponder(http.StatusOK, expectedResult.Text)
		assert.NoError(t, err)
		responderBinary, err := httpmock.NewJsonResponder(http.StatusOK, expectedResult.Binary)
		assert.NoError(t, err)
		responderCard, err := httpmock.NewJsonResponder(http.StatusOK, expectedResult.Card)
		assert.NoError(t, err)
		responderCredential, err := httpmock.NewJsonResponder(
			http.StatusOK,
			expectedResult.Credential,
		)
		assert.NoError(t, err)

		httpmock.RegisterResponder(
			http.MethodGet,
			fmt.Sprintf("%v/api/store/%v", baseURL, srvrModels.TextCollection),
			responderText,
		)
		httpmock.RegisterResponder(
			http.MethodGet,
			fmt.Sprintf("%v/api/store/%v", baseURL, srvrModels.BinaryCollection),
			responderBinary,
		)
		httpmock.RegisterResponder(
			http.MethodGet,
			fmt.Sprintf("%v/api/store/%v", baseURL, srvrModels.CardCollection),
			responderCard,
		)
		httpmock.RegisterResponder(
			http.MethodGet,
			fmt.Sprintf("%v/api/store/%v", baseURL, srvrModels.CredentialsCollection),
			responderCredential,
		)

		// Invoke Sync method
		actualResult, actualError := s.Sync("good-token", collectionNames)
		// Assert results
		assert.Equal(t, expectedResult, actualResult)
		assert.NoError(t, actualError)
	})

	t.Run("unauthorized", func(t *testing.T) {
		httpmock.Reset()
		collectionNames := []srvrModels.CollectionName{
			srvrModels.TextCollection,
		}
		responder := httpmock.NewStringResponder(http.StatusUnauthorized, "")
		httpmock.RegisterResponder(
			http.MethodGet,
			fmt.Sprintf("%v/api/store/%v", baseURL, srvrModels.TextCollection),
			responder,
		)

		// Invoke Sync method
		_, actualError := s.Sync("bad-token", collectionNames)
		// Assert results
		assert.ErrorIs(t, srvErrors.ErrUnauthorized, actualError)
	})
	t.Run("bad_request", func(t *testing.T) {
		httpmock.Reset()
		collectionNames := []srvrModels.CollectionName{
			srvrModels.TextCollection,
		}
		responder := httpmock.NewStringResponder(http.StatusUnauthorized, "")
		httpmock.RegisterResponder(
			http.MethodGet,
			fmt.Sprintf("%v/api/store/%v", baseURL, srvrModels.TextCollection),
			responder,
		)

		expectedResult := &clientModels.SyncResponse{
			Text: []srvrModels.TextRecord{
				{RecordID: primitive.NewObjectID(), Data: srvrModels.TextInfo("text1")},
			},
		}

		responder, err := httpmock.NewJsonResponder(http.StatusOK, expectedResult)
		assert.NoError(t, err)
		httpmock.RegisterResponder(
			http.MethodGet,
			fmt.Sprintf("%v/api/store/%v", baseURL, srvrModels.TextCollection),
			responder,
		)
		// Invoke Sync method
		_, actualError := s.Sync("bad-token", collectionNames)
		// Assert results
		assert.Error(t, actualError)
	})

}

func TestSyncService_Register(t *testing.T) {
	baseURL := "https://example.com"
	s := NewSyncService(baseURL)
	client := s.GetClient()
	assert.Equal(t, baseURL, client.HostURL)

	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	t.Run("success", func(t *testing.T) {
		httpmock.Reset()

		httpmock.RegisterResponder(
			http.MethodPost,
			fmt.Sprintf("%v/api/sync/register", baseURL),
			httpmock.NewStringResponder(http.StatusOK, "ok"),
		)

		r, err := s.Register("some-token", "http://localhost:1234")

		assert.Equal(t, "ok", r)
		assert.NoError(t, err)

	})
	t.Run("fail", func(t *testing.T) {
		httpmock.Reset()

		httpmock.RegisterResponder(
			http.MethodPost,
			fmt.Sprintf("%v/api/sync/register", baseURL),
			httpmock.NewStringResponder(http.StatusBadRequest, "some error"),
		)

		r, err := s.Register("some-token", "http://localhost:1234")

		assert.Equal(t, "", r)
		assert.Error(t, err)

	})
}

func TestSyncService_Unregister(t *testing.T) {
	baseURL := "https://example.com"
	s := NewSyncService(baseURL)
	client := s.GetClient()
	assert.Equal(t, baseURL, client.HostURL)

	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	t.Run("success", func(t *testing.T) {
		httpmock.Reset()

		httpmock.RegisterResponder(
			http.MethodPost,
			fmt.Sprintf("%v/api/sync/unregister", baseURL),
			httpmock.NewStringResponder(http.StatusOK, "ok"),
		)

		r, err := s.Unregister("some-token", "http://localhost:1234")

		assert.Equal(t, "ok", r)
		assert.NoError(t, err)

	})
	t.Run("fail", func(t *testing.T) {
		httpmock.Reset()

		httpmock.RegisterResponder(
			http.MethodPost,
			fmt.Sprintf("%v/api/sync/unregister", baseURL),
			httpmock.NewStringResponder(http.StatusBadRequest, "some error"),
		)

		r, err := s.Unregister("some-token", "http://localhost:1234")

		assert.Equal(t, "", r)
		assert.Error(t, err)

	})
}
