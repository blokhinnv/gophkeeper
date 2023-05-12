package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestUtilsService(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("constructor", func(mt *mtest.T) {
		s := NewUtilsService(mt.Client)
		assert.NotNil(t, s, "NewUtilsService should return a non-nil pointer.")
	})
	mt.Run("success", func(mt *mtest.T) {
		s := NewUtilsService(mt.Client)
		err := s.Ping()
		require.Error(t, err)
	})
}
