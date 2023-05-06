package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	"github.com/blokhinnv/gophkeeper/internal/server/models"
	"github.com/blokhinnv/gophkeeper/pkg/encrypt"
)

type StorageServiceTestSuite struct {
	suite.Suite
}

func (suite *StorageServiceTestSuite) SetupSuite()    {}
func (suite *StorageServiceTestSuite) TearDownSuite() {}

func (suite *StorageServiceTestSuite) TestStore() {
	t := suite.T()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("success_text", func(mt *mtest.T) {
		storageService := NewStorageService(mt.DB, "my-secret-key")
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		rec := models.UntypedRecord{
			UntypedRecordContent: models.UntypedRecordContent{
				Data: "test message",
			},
			Username: "blokhinnv",
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		res, err := storageService.Store(context.TODO(), models.TextCollection, rec)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
	mt.Run("success_not_text", func(mt *mtest.T) {
		storageService := NewStorageService(mt.DB, "my-secret-key")
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		rec := models.UntypedRecord{
			UntypedRecordContent: models.UntypedRecordContent{
				Data: map[string]any{
					"login":    "blokhinnv",
					"password": "test-password",
				},
			},
			Username: "blokhinnv",
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		res, err := storageService.Store(context.TODO(), models.CredentialsCollection, rec)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}

func (suite *StorageServiceTestSuite) TestGetAll() {
	t := suite.T()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("success_text", func(mt *mtest.T) {
		secretKey := "my-secret-key"
		storageService := NewStorageService(mt.DB, secretKey)

		username := "blokhinnv"
		rawData := "some text data.."
		data, err := encrypt.EncryptString(rawData, secretKey)
		require.NoError(t, err)

		batchItem := mtest.CreateCursorResponse(1, "get_all.success_text", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "usename", Value: username},
			{Key: "data", Value: data},
		})
		batchEnd := mtest.CreateCursorResponse(0, "get_all.success_text", mtest.NextBatch)
		mt.AddMockResponses(batchItem, batchEnd)

		res, err := storageService.GetAll(context.TODO(), models.TextCollection, username)
		require.NoError(t, err)
		require.NotEmpty(t, res)
		require.Equal(t, rawData, res[0].Data)
	})
	mt.Run("success_not_text", func(mt *mtest.T) {
		secretKey := "my-secret-key"
		storageService := NewStorageService(mt.DB, secretKey)

		username := "blokhinnv"
		rawData := map[string]any{
			"login":    "some-user",
			"password": "some-password",
		}
		data, err := encrypt.EncryptMap(rawData, secretKey)
		require.NoError(t, err)

		batchItem := mtest.CreateCursorResponse(
			1,
			"get_all.success_not_text",
			mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: primitive.NewObjectID()},
				{Key: "usename", Value: username},
				{Key: "data", Value: data},
			},
		)
		batchEnd := mtest.CreateCursorResponse(0, "get_all.success_not_text", mtest.NextBatch)
		mt.AddMockResponses(batchItem, batchEnd)

		res, err := storageService.GetAll(context.TODO(), models.TextCollection, username)
		require.NoError(t, err)
		require.NotEmpty(t, res)

		resMap, ok := res[0].Data.(map[string]any)
		require.Equal(t, true, ok)
		require.Equal(t, rawData["login"], resMap["login"])
		require.Equal(t, rawData["password"], resMap["password"])
	})
	mt.Run("empty_response", func(mt *mtest.T) {
		secretKey := "my-secret-key"
		storageService := NewStorageService(mt.DB, secretKey)

		username := "blokhinnv"
		res, err := storageService.GetAll(context.TODO(), models.TextCollection, username)
		require.Error(t, err)
		require.Empty(t, res)
	})
}

func (suite *StorageServiceTestSuite) TestUpdate() {
	t := suite.T()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("success_text", func(mt *mtest.T) {
		storageService := NewStorageService(mt.DB, "my-secret-key")
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "nModified", Value: 1},
		})

		err := storageService.Update(
			context.TODO(),
			models.TextCollection,
			"blokhinnv",
			primitive.NewObjectID(),
			"test message",
			make(models.Metadata),
		)
		require.NoError(t, err)
	})
	mt.Run("success_not_text", func(mt *mtest.T) {
		storageService := NewStorageService(mt.DB, "my-secret-key")
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "nModified", Value: 1},
		})

		err := storageService.Update(
			context.TODO(),
			models.CredentialsCollection,
			"blokhinnv",
			primitive.NewObjectID(),
			map[string]any{"login": "blokhinnv", "password": "some-pwd"},
			make(models.Metadata),
		)
		require.NoError(t, err)
	})
	mt.Run("error", func(mt *mtest.T) {
		storageService := NewStorageService(mt.DB, "my-secret-key")
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 0},
		})

		err := storageService.Update(
			context.TODO(),
			models.CredentialsCollection,
			"blokhinnv",
			primitive.NewObjectID(),
			map[string]any{"login": "blokhinnv", "password": "some-pwd"},
			make(models.Metadata),
		)
		require.Error(t, err)
	})
}

func (suite *StorageServiceTestSuite) TestDelete() {
	t := suite.T()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("success", func(mt *mtest.T) {
		storageService := NewStorageService(mt.DB, "my-secret-key")
		mt.AddMockResponses(
			bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 1}},
		)
		err := storageService.Delete(
			context.TODO(),
			models.TextCollection,
			"blokhinnv",
			primitive.NewObjectID(),
		)
		require.NoError(t, err)
	})
	mt.Run("error", func(mt *mtest.T) {
		storageService := NewStorageService(mt.DB, "my-secret-key")
		mt.AddMockResponses(
			bson.D{{Key: "ok", Value: 0}},
		)
		err := storageService.Delete(
			context.TODO(),
			models.TextCollection,
			"blokhinnv",
			primitive.NewObjectID(),
		)
		require.Error(t, err)
	})

}

func TestStorageServiceTestSuite(t *testing.T) {
	suite.Run(t, new(StorageServiceTestSuite))
}
