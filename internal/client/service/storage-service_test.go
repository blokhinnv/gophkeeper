package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	clientModels "github.com/blokhinnv/gophkeeper/internal/client/models"
	"github.com/blokhinnv/gophkeeper/internal/server/models"
	srvrModels "github.com/blokhinnv/gophkeeper/internal/server/models"
)

func TestStorageService_GetAll(t *testing.T) {
	baseURL := "https://example.com"
	s := NewStorageService(baseURL)
	data := &clientModels.SyncResponse{
		Text: []srvrModels.TextRecord{
			{
				RecordID: models.NewRandomObjectID(),
				Username: "blokhinnv",
				Data:     srvrModels.TextInfo("some data..."),
			},
		},
	}

	t.Run("get_text", func(t *testing.T) {
		r := s.GetAll(srvrModels.TextCollection, data)
		// Assert
		texts, ok := r.([]srvrModels.TextRecord)
		if !ok {
			t.Errorf("Expected result to be of type []models.Text, but got %T", r)
		}
		if len(texts) != len(data.Text) {
			t.Errorf("Expected %d texts, but got %d", len(data.Text), len(texts))
		}
		for i, text := range data.Text {
			if texts[i].RecordID != text.RecordID || texts[i].Data != text.Data {
				t.Errorf("Expected text %d to be %+v, but got %+v", i, text, texts[i])
			}
		}
	})
	t.Run("get_creds", func(t *testing.T) {
		r := s.GetAll(srvrModels.CredentialsCollection, data)
		assert.Nil(t, r)
	})
	t.Run("other", func(t *testing.T) {
		r := s.GetAll(srvrModels.CollectionName("other"), data)
		assert.Nil(t, r)
	})
}

func TestStorageService_Add(t *testing.T) {
	baseURL := "https://example.com"
	s := NewStorageService(baseURL)
	client := s.GetClient()
	assert.Equal(t, baseURL, client.HostURL)

	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	t.Run("ok", func(t *testing.T) {
		httpmock.Reset()

		httpmock.RegisterResponder(
			http.MethodPut,
			fmt.Sprintf("%v/api/store/%v", baseURL, srvrModels.TextCollection),
			httpmock.NewStringResponder(200, "ok"),
		)

		body := srvrModels.TextRecord{
			RecordID: models.NewRandomObjectID(),
			Username: "blokhinnv",
			Data:     srvrModels.TextInfo("some data..."),
		}

		bodyEncoded, err := json.Marshal(body)
		assert.NoError(t, err)

		resp, err := s.Add(string(bodyEncoded), srvrModels.TextCollection, "some-token...")
		assert.NoError(t, err)
		assert.Equal(t, "ok", resp)
	})
	t.Run("bad", func(t *testing.T) {
		httpmock.Reset()

		httpmock.RegisterResponder(
			http.MethodPut,
			fmt.Sprintf("%v/api/store/%v", baseURL, srvrModels.TextCollection),
			httpmock.NewStringResponder(400, "bad"),
		)

		bodyEncoded := `{"bad"}`

		resp, err := s.Add(string(bodyEncoded), srvrModels.TextCollection, "some-token...")
		assert.Equal(t, "", resp)
		assert.Equal(t, "bad", err.Error())
	})
}

func TestStorageService_Update(t *testing.T) {
	baseURL := "https://example.com"
	s := NewStorageService(baseURL)
	client := s.GetClient()
	assert.Equal(t, baseURL, client.HostURL)

	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	t.Run("ok", func(t *testing.T) {
		httpmock.Reset()

		httpmock.RegisterResponder(
			http.MethodPost,
			fmt.Sprintf("%v/api/store/%v", baseURL, srvrModels.TextCollection),
			httpmock.NewStringResponder(200, "ok"),
		)

		body := srvrModels.TextRecord{
			RecordID: models.NewRandomObjectID(),
			Username: "blokhinnv",
			Data:     srvrModels.TextInfo("some data..."),
		}

		bodyEncoded, err := json.Marshal(body)
		assert.NoError(t, err)

		resp, err := s.Update(string(bodyEncoded), srvrModels.TextCollection, "some-token...")
		assert.NoError(t, err)
		assert.Equal(t, "ok", resp)
	})
	t.Run("bad", func(t *testing.T) {
		httpmock.Reset()

		httpmock.RegisterResponder(
			http.MethodPost,
			fmt.Sprintf("%v/api/store/%v", baseURL, srvrModels.TextCollection),
			httpmock.NewStringResponder(400, "bad"),
		)

		bodyEncoded := `{"bad"}`

		resp, err := s.Update(string(bodyEncoded), srvrModels.TextCollection, "some-token...")
		assert.Equal(t, "", resp)
		assert.Equal(t, "bad", err.Error())
	})
}
func TestStorageService_Delete(t *testing.T) {
	baseURL := "https://example.com"
	s := NewStorageService(baseURL)
	client := s.GetClient()
	assert.Equal(t, baseURL, client.HostURL)

	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	t.Run("ok", func(t *testing.T) {
		httpmock.Reset()

		httpmock.RegisterResponder(
			http.MethodDelete,
			fmt.Sprintf("%v/api/store/%v", baseURL, srvrModels.TextCollection),
			httpmock.NewStringResponder(200, "ok"),
		)

		body := fmt.Sprintf(`{"record_id": "%v"}`, models.NewRandomObjectID())

		resp, err := s.Delete(body, srvrModels.TextCollection, "some-token...")
		assert.NoError(t, err)
		assert.Equal(t, "ok", resp)
	})
	t.Run("bad", func(t *testing.T) {
		httpmock.Reset()

		httpmock.RegisterResponder(
			http.MethodDelete,
			fmt.Sprintf("%v/api/store/%v", baseURL, srvrModels.TextCollection),
			httpmock.NewStringResponder(400, "bad"),
		)

		body := fmt.Sprintf(`{"qwe": "%v"}`, models.NewRandomObjectID())
		resp, err := s.Delete(body, srvrModels.TextCollection, "some-token...")
		assert.Equal(t, "", resp)
		assert.Equal(t, "bad", err.Error())
	})
}
