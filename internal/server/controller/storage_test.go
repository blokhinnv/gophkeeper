package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	srvErrors "github.com/blokhinnv/gophkeeper/internal/server/errors"
	"github.com/blokhinnv/gophkeeper/internal/server/middleware"
	"github.com/blokhinnv/gophkeeper/internal/server/models"
	"github.com/blokhinnv/gophkeeper/internal/server/service/mock"
)

func TestNewStorageController(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// create a new storageController instance with a mocked service
	storage := mock.NewMockStorageService(mockCtrl)
	sync := mock.NewMockSyncService(mockCtrl)
	ctrl := NewStorageController(storage, sync)
	assert.NotNil(t, ctrl)
}

func TestStorageController_ValidateDataField(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// create a new storageController instance with a mocked service
	storage := mock.NewMockStorageService(mockCtrl)
	sync := mock.NewMockSyncService(mockCtrl)
	ctrl, ok := NewStorageController(storage, sync).(*storageController)
	assert.NotNil(t, ctrl)
	assert.Equal(t, true, ok)

	t.Run("ok_credentials", func(t *testing.T) {
		// Test case 1: Successful validation for CredentialsCollection
		credentials := models.CredentialInfo{
			Login:    "test",
			Password: "password",
		}
		err := ctrl.validateDataField(credentials, models.CredentialsCollection)
		assert.NoError(t, err)
	})
	t.Run("bad_credentials", func(t *testing.T) {
		// Test case 2: Failed validation for CredentialsCollection
		credentials := models.CredentialInfo{
			Password: "password",
		}
		err := ctrl.validateDataField(credentials, models.CredentialsCollection)
		assert.Error(t, err)
	})
	t.Run("ok_card", func(t *testing.T) {
		// Test case 3: Successful validation for CardCollection
		card := models.CardInfo{
			CardNumber:     "4111 1111 1111 1111",
			ExpirationDate: "01/12",
			CVV:            "123",
		}
		err := ctrl.validateDataField(card, models.CardCollection)
		assert.NoError(t, err)
	})
	t.Run("bad_card", func(t *testing.T) {
		// Test case 4: Failed validation for CardCollection
		card := models.CardInfo{
			CardNumber:     "1234",
			ExpirationDate: "01/12",
			CVV:            "123",
		}
		err := ctrl.validateDataField(card, models.CardCollection)
		assert.Error(t, err)
	})
	t.Run("ok_binary", func(t *testing.T) {
		// Test case 5: Successful validation for BinaryCollection
		binary := models.BinaryInfo{
			FileName: "file.txt",
			Content:  "cXdlcg==",
		}
		err := ctrl.validateDataField(binary, models.BinaryCollection)
		assert.NoError(t, err)
	})
	t.Run("bad_binary", func(t *testing.T) {
		// Test case 6: Failed validation for BinaryCollection
		binary := models.BinaryInfo{
			FileName: "file.txt",
			Content:  "cXg=123",
		}
		err := ctrl.validateDataField(binary, models.BinaryCollection)
		assert.Error(t, err)
	})
	t.Run("other", func(t *testing.T) {
		// Test case 7: No validation required for other collections
		data := "test"
		err := ctrl.validateDataField(data, "text")
		assert.NoError(t, err)
	})
}

func TestStorageController_Store(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// create a new storageController instance with a mocked service
	storage := mock.NewMockStorageService(mockCtrl)
	sync := mock.NewMockSyncService(mockCtrl)
	ctrl := NewStorageController(storage, sync)
	assert.NotNil(t, ctrl)

	t.Run("no_username", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/", nil)
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req

		ctrl.Store(ctx)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, srvErrors.ErrNoUsernameProvided.Error(), rec.Body.String())
	})
	t.Run("bad_collection", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/collections/invalid-collection", nil)
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req
		ctx.Set(middleware.UsernameContextValue, "username")
		ctx.Params = append(
			ctx.Params,
			gin.Param{Key: "collectionName", Value: "invalid-collection"},
		)

		ctrl.Store(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("bad_body", func(t *testing.T) {
		reqBody := bytes.NewBufferString(`{"data": "invalid-json"}`)
		req, _ := http.NewRequest("POST", "/collections/credentials", reqBody)
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req
		ctx.Params = append(ctx.Params, gin.Param{Key: "collectionName", Value: "credentials"})
		ctx.Set(middleware.UsernameContextValue, "username")

		ctrl.Store(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("not_valid_body", func(t *testing.T) {
		reqBody := bytes.NewBufferString(`{"data": {"login": ""}}`)
		req, _ := http.NewRequest("POST", "/collections/credentials", reqBody)
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req
		ctx.Params = append(ctx.Params, gin.Param{Key: "collectionName", Value: "credentials"})
		ctx.Set(middleware.UsernameContextValue, "username")

		ctrl.Store(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("service_err", func(t *testing.T) {
		storage.EXPECT().
			Store(gomock.Any(), gomock.Eq(models.CollectionName("credentials")), gomock.Any()).
			Return("", fmt.Errorf("some error"))
		reqBody := bytes.NewBufferString(
			`{"data": {"login": "user123", "password": "password123"}}`,
		)
		req, _ := http.NewRequest("POST", "/collections/credentials", reqBody)
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req
		ctx.Params = append(ctx.Params, gin.Param{Key: "collectionName", Value: "credentials"})
		ctx.Set(middleware.UsernameContextValue, "username")

		ctrl.Store(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("ok", func(t *testing.T) {
		storage.EXPECT().
			Store(gomock.Any(), gomock.Eq(models.CollectionName("credentials")), gomock.Any()).
			Return("some-id", nil)
		sync.EXPECT().Signal(gomock.Any()).AnyTimes()
		reqBody := bytes.NewBufferString(
			`{"data": {"login": "user123", "password": "password123"}}`,
		)
		req, _ := http.NewRequest("POST", "/collections/credentials", reqBody)
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req
		ctx.Params = append(ctx.Params, gin.Param{Key: "collectionName", Value: "credentials"})
		ctx.Set(middleware.UsernameContextValue, "username")

		ctrl.Store(ctx)

		assert.Equal(t, http.StatusAccepted, rec.Code)
		assert.Contains(
			t,
			rec.Body.String(),
			"Record added to credentials collection",
		)
	})
}

func TestStorageController_GetAll(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// create a new storageController instance with a mocked service
	storage := mock.NewMockStorageService(mockCtrl)
	sync := mock.NewMockSyncService(mockCtrl)
	ctrl := NewStorageController(storage, sync)
	assert.NotNil(t, ctrl)

	username := "testuser"
	t.Run("ok", func(t *testing.T) {
		expectedRecords := []models.UntypedRecord{
			{
				RecordID: models.NewRandomObjectID(),
				Username: username,
				UntypedRecordContent: models.UntypedRecordContent{
					Data: map[string]any{
						"login":    "john",
						"password": "password123",
					},
					Metadata: models.Metadata{
						"created_at": "2023-05-06",
					},
				},
			},
		}
		storage.EXPECT().
			GetAll(gomock.Any(), models.CollectionName("credentials"), username).
			Return(expectedRecords, nil)
		req, _ := http.NewRequest("GET", "/collections/credentials", nil)
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req
		ctx.Params = append(ctx.Params, gin.Param{Key: "collectionName", Value: "credentials"})
		ctx.Set(middleware.UsernameContextValue, username)

		ctrl.GetAll(ctx)

		assert.Equal(t, http.StatusOK, rec.Code)
		var responseRecords []models.UntypedRecord
		err := json.Unmarshal(rec.Body.Bytes(), &responseRecords)
		assert.NoError(t, err)
		// Do not return username
		expectedRecords[0].Username = ""
		assert.Equal(t, expectedRecords, responseRecords)

	})

	t.Run("bad_response", func(t *testing.T) {
		storage.EXPECT().
			GetAll(gomock.Any(), models.CollectionName("credentials"), username).
			Return(nil, fmt.Errorf("some error"))
		req, _ := http.NewRequest("GET", "/collections/credentials", nil)
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req
		ctx.Params = append(ctx.Params, gin.Param{Key: "collectionName", Value: "credentials"})
		ctx.Set(middleware.UsernameContextValue, username)

		ctrl.GetAll(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("no_username", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/collections/text", nil)
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req

		ctrl.GetAll(ctx)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, srvErrors.ErrNoUsernameProvided.Error(), rec.Body.String())
	})

	t.Run("bad_collection", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/collections/invalid-collection", nil)
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req
		ctx.Set(middleware.UsernameContextValue, username)
		ctx.Params = append(
			ctx.Params,
			gin.Param{Key: "collectionName", Value: "invalid-collection"},
		)

		ctrl.GetAll(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

}

func TestStorageController_Update(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// create a new storageController instance with a mocked service
	storage := mock.NewMockStorageService(mockCtrl)
	sync := mock.NewMockSyncService(mockCtrl)
	ctrl := NewStorageController(storage, sync)
	assert.NotNil(t, ctrl)
	username := "testuser"

	t.Run("no_username", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/store/credentials", nil)
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req

		ctrl.Update(ctx)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, srvErrors.ErrNoUsernameProvided.Error(), rec.Body.String())
	})

	t.Run("bad_collection", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/store/invalid-collection", nil)
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req
		ctx.Set(middleware.UsernameContextValue, username)
		ctx.Params = append(
			ctx.Params,
			gin.Param{Key: "collectionName", Value: "invalid-collection"},
		)

		ctrl.Update(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ok", func(t *testing.T) {
		storage.EXPECT().
			Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)
		sync.EXPECT().Signal(gomock.Any()).AnyTimes()
		recordID := models.NewRandomObjectID()
		record := models.UntypedRecord{
			RecordID: recordID,
			Username: username,
			UntypedRecordContent: models.UntypedRecordContent{
				Data: map[string]any{
					"login":    "john",
					"password": "password123",
				},
				Metadata: models.Metadata{
					"created_at": "2023-05-06",
				},
			},
		}

		data, _ := json.Marshal(record)
		req, _ := http.NewRequest("POST", "/api/store/credentials", bytes.NewBuffer(data))
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req
		ctx.Set(middleware.UsernameContextValue, username)
		ctx.Params = append(
			ctx.Params,
			gin.Param{Key: "collectionName", Value: "credentials"},
		)

		ctrl.Update(ctx)

		assert.Equal(t, http.StatusAccepted, rec.Code)
	})
	t.Run("not_found", func(t *testing.T) {
		storage.EXPECT().
			Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(srvErrors.ErrRecordNotFound)
		sync.EXPECT().Signal(gomock.Any()).AnyTimes()
		recordID := models.NewRandomObjectID()
		record := models.UntypedRecord{
			RecordID: recordID,
			Username: username,
			UntypedRecordContent: models.UntypedRecordContent{
				Data: map[string]any{
					"login":    "john",
					"password": "password123",
				},
				Metadata: models.Metadata{
					"created_at": "2023-05-06",
				},
			},
		}

		data, _ := json.Marshal(record)
		req, _ := http.NewRequest("POST", "/api/store/credentials", bytes.NewBuffer(data))
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req
		ctx.Set(middleware.UsernameContextValue, username)
		ctx.Params = append(
			ctx.Params,
			gin.Param{Key: "collectionName", Value: "credentials"},
		)

		ctrl.Update(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("service_error", func(t *testing.T) {
		storage.EXPECT().
			Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(fmt.Errorf("some error"))
		sync.EXPECT().Signal(gomock.Any()).AnyTimes()
		recordID := models.NewRandomObjectID()
		record := models.UntypedRecord{
			RecordID: recordID,
			Username: username,
			UntypedRecordContent: models.UntypedRecordContent{
				Data: map[string]any{
					"login":    "john",
					"password": "password123",
				},
				Metadata: models.Metadata{
					"created_at": "2023-05-06",
				},
			},
		}

		data, _ := json.Marshal(record)
		req, _ := http.NewRequest("POST", "/api/store/credentials", bytes.NewBuffer(data))
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req
		ctx.Set(middleware.UsernameContextValue, username)
		ctx.Params = append(
			ctx.Params,
			gin.Param{Key: "collectionName", Value: "credentials"},
		)

		ctrl.Update(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
	t.Run("not_valid_data", func(t *testing.T) {
		recordID := models.NewRandomObjectID()
		record := models.UntypedRecord{
			RecordID: recordID,
			Username: username,
			UntypedRecordContent: models.UntypedRecordContent{
				Data: map[string]any{
					"login":    "",
					"password": "",
				},
				Metadata: models.Metadata{
					"created_at": "2023-05-06",
				},
			},
		}

		data, _ := json.Marshal(record)
		req, _ := http.NewRequest("POST", "/api/store/credentials", bytes.NewBuffer(data))
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req
		ctx.Set(middleware.UsernameContextValue, username)
		ctx.Params = append(
			ctx.Params,
			gin.Param{Key: "collectionName", Value: "credentials"},
		)

		ctrl.Update(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("bad_data", func(t *testing.T) {
		recordID := models.NewRandomObjectID()
		record := models.UntypedRecord{
			RecordID: recordID,
			Username: username,
			UntypedRecordContent: models.UntypedRecordContent{
				Data: map[string]any{
					"login":    543,
					"password": 123,
				},
				Metadata: models.Metadata{
					"created_at": "2023-05-06",
				},
			},
		}

		data, _ := json.Marshal(record)
		req, _ := http.NewRequest("POST", "/api/store/credentials", bytes.NewBuffer(data))
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req
		ctx.Set(middleware.UsernameContextValue, username)
		ctx.Params = append(
			ctx.Params,
			gin.Param{Key: "collectionName", Value: "credentials"},
		)

		ctrl.Update(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
func TestStorageController_Delete(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// create a new storageController instance with a mocked service
	storage := mock.NewMockStorageService(mockCtrl)
	sync := mock.NewMockSyncService(mockCtrl)
	ctrl := NewStorageController(storage, sync)
	assert.NotNil(t, ctrl)
	username := "testuser"

	type deleteRequestBody struct {
		RecordID models.ObjectID `json:"record_id" bson:"_id" binding:"required"`
	}

	t.Run("no_username", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/store/credentials", nil)
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req

		ctrl.Delete(ctx)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, srvErrors.ErrNoUsernameProvided.Error(), rec.Body.String())
	})

	t.Run("bad_collection", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/store/invalid-collection", nil)
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req
		ctx.Set(middleware.UsernameContextValue, username)
		ctx.Params = append(
			ctx.Params,
			gin.Param{Key: "collectionName", Value: "invalid-collection"},
		)

		ctrl.Delete(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ok", func(t *testing.T) {
		storage.EXPECT().
			Delete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)
		sync.EXPECT().Signal(gomock.Any()).AnyTimes()

		recordID := models.NewRandomObjectID()
		record := deleteRequestBody{recordID}
		data, _ := json.Marshal(record)

		req, _ := http.NewRequest("DELETE", "/api/store/credentials", bytes.NewBuffer(data))
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req
		ctx.Set(middleware.UsernameContextValue, username)
		ctx.Params = append(
			ctx.Params,
			gin.Param{Key: "collectionName", Value: "credentials"},
		)

		ctrl.Delete(ctx)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
	t.Run("service_error", func(t *testing.T) {
		storage.EXPECT().
			Delete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(fmt.Errorf("some error"))
		sync.EXPECT().Signal(gomock.Any()).AnyTimes()

		recordID := models.NewRandomObjectID()
		record := deleteRequestBody{recordID}
		data, _ := json.Marshal(record)

		req, _ := http.NewRequest("DELETE", "/api/store/credentials", bytes.NewBuffer(data))
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req
		ctx.Set(middleware.UsernameContextValue, username)
		ctx.Params = append(
			ctx.Params,
			gin.Param{Key: "collectionName", Value: "credentials"},
		)

		ctrl.Delete(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
	t.Run("not_found", func(t *testing.T) {
		storage.EXPECT().
			Delete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(srvErrors.ErrRecordNotFound)
		sync.EXPECT().Signal(gomock.Any()).AnyTimes()

		recordID := models.NewRandomObjectID()
		record := deleteRequestBody{recordID}
		data, _ := json.Marshal(record)

		req, _ := http.NewRequest("DELETE", "/api/store/credentials", bytes.NewBuffer(data))
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req
		ctx.Set(middleware.UsernameContextValue, username)
		ctx.Params = append(
			ctx.Params,
			gin.Param{Key: "collectionName", Value: "credentials"},
		)

		ctrl.Delete(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
