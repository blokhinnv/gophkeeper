package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceTestSuite struct {
	suite.Suite
}

func (suite *AuthServiceTestSuite) SetupSuite()    {}
func (suite *AuthServiceTestSuite) TearDownSuite() {}

func (suite *AuthServiceTestSuite) TestRegister() {
	t := suite.T()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("success", func(mt *mtest.T) {
		// Test registering a new user
		authService := NewAuthService(mt.Coll, "my-secret-key", time.Hour)
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		err := authService.Register("testuser", "testpassword")
		require.NoError(t, err)
	})
	mt.Run("err", func(mt *mtest.T) {
		// Test registering a user with an already taken username
		authService := NewAuthService(mt.Coll, "my-secret-key", time.Hour)
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))
		err := authService.Register("testuser", "testpassword2")
		require.Error(t, err)
	})
}

func (suite *AuthServiceTestSuite) TestLogin() {
	t := suite.T()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	user := "testuser"
	password := "testpassword"
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	require.NoError(t, err)

	mt.Run("valid credentials", func(mt *mtest.T) {
		authService := NewAuthService(mt.Coll, "my-secret-key", time.Hour)
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "login.valid", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "username", Value: user},
			{Key: "hashedPassword", Value: string(hashedPassword)},
		}))
		tok, err := authService.Login(user, password)
		require.NoError(t, err)
		require.NotEmpty(t, tok)
	})
	mt.Run("invalid credentials", func(mt *mtest.T) {
		authService := NewAuthService(mt.Coll, "my-secret-key", time.Hour)
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "login.valid", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "username", Value: user},
			{Key: "hashedPassword", Value: string(hashedPassword)},
		}))
		_, err := authService.Login(user, "test-password-2")
		require.Error(t, err)
	})
	mt.Run("invalid credentials", func(mt *mtest.T) {
		authService := NewAuthService(mt.Coll, "my-secret-key", time.Hour)
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Message: "mongo: no documents in result",
		}))
		_, err := authService.Login(user, password)
		require.Error(t, err)
	})
}

func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}
