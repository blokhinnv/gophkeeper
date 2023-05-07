package service

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	srvErrors "github.com/blokhinnv/gophkeeper/internal/server/errors"
	"github.com/blokhinnv/gophkeeper/internal/server/models"
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
		collectionNames := []models.CollectionName{
			models.TextCollection,
			models.BinaryCollection,
			models.CardCollection,
			models.CredentialsCollection,
		}
		expectedResult := &syncResponse{
			Text: []models.TextRecord{
				{RecordID: primitive.NewObjectID(), Data: models.TextInfo("some-text...")},
			},
			Binary: []models.BinaryRecord{
				{
					RecordID: primitive.NewObjectID(),
					Data:     models.BinaryInfo{FileName: "test.test", Content: "cXdlcXdld3Fl"},
				},
			},
			Card: []models.CardRecord{
				{
					RecordID: primitive.NewObjectID(),
					Data: models.CardInfo{
						CardNumber:     "1234 1234 1234 1234",
						CVV:            "234",
						ExpirationDate: "01/12",
					},
				},
			},
			Credential: []models.CredentialRecord{
				{
					RecordID: primitive.NewObjectID(),
					Data:     models.CredentialInfo{Login: "some-login", Password: "some-password"},
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
			fmt.Sprintf("%v/api/store/%v", baseURL, models.TextCollection),
			responderText,
		)
		httpmock.RegisterResponder(
			http.MethodGet,
			fmt.Sprintf("%v/api/store/%v", baseURL, models.BinaryCollection),
			responderBinary,
		)
		httpmock.RegisterResponder(
			http.MethodGet,
			fmt.Sprintf("%v/api/store/%v", baseURL, models.CardCollection),
			responderCard,
		)
		httpmock.RegisterResponder(
			http.MethodGet,
			fmt.Sprintf("%v/api/store/%v", baseURL, models.CredentialsCollection),
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
		collectionNames := []models.CollectionName{
			models.TextCollection,
		}
		responder := httpmock.NewStringResponder(http.StatusUnauthorized, "")
		httpmock.RegisterResponder(
			http.MethodGet,
			fmt.Sprintf("%v/api/store/%v", baseURL, models.TextCollection),
			responder,
		)

		// Invoke Sync method
		_, actualError := s.Sync("bad-token", collectionNames)
		// Assert results
		assert.ErrorIs(t, srvErrors.ErrUnauthorized, actualError)
	})
	t.Run("bad_request", func(t *testing.T) {
		httpmock.Reset()
		collectionNames := []models.CollectionName{
			models.TextCollection,
		}
		responder := httpmock.NewStringResponder(http.StatusUnauthorized, "")
		httpmock.RegisterResponder(
			http.MethodGet,
			fmt.Sprintf("%v/api/store/%v", baseURL, models.TextCollection),
			responder,
		)

		expectedResult := &syncResponse{
			Text: []models.TextRecord{
				{RecordID: primitive.NewObjectID(), Data: models.TextInfo("text1")},
			},
		}

		responder, err := httpmock.NewJsonResponder(http.StatusOK, expectedResult)
		assert.NoError(t, err)
		httpmock.RegisterResponder(
			http.MethodGet,
			fmt.Sprintf("%v/api/store/%v", baseURL, models.TextCollection),
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
