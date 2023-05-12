package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/blokhinnv/gophkeeper/internal/server/middleware"
	"github.com/blokhinnv/gophkeeper/internal/server/service/mock"
)

func TestSyncController_Register(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	sync := mock.NewMockSyncService(mockCtrl)
	ctrl := NewSyncController(sync)
	assert.NotNil(t, ctrl)

	t.Run("no_username", func(t *testing.T) {
		// Test case 1: Missing username
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.POST("/register", ctrl.Register)
		req, _ := http.NewRequest(http.MethodPost, "/register", nil)
		c.Request = req
		ctrl.Register(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
	t.Run("invalid_data", func(t *testing.T) {
		// Test case 2: Invalid JSON payload
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		c.Set(middleware.UsernameContextValue, "blokhinnv")
		r.POST("/register", ctrl.Register)
		req, _ := http.NewRequest(
			http.MethodPost,
			"/register",
			bytes.NewBufferString(`{"invalid": json}`),
		)
		c.Request = req
		ctrl.Register(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ok", func(t *testing.T) {
		// Test case 3: Valid registration
		sync.EXPECT().Register(gomock.Any()).Times(1)
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		c.Set(middleware.UsernameContextValue, "blokhinnv")
		r.POST("/register", ctrl.Register)
		req, _ := http.NewRequest(
			http.MethodPost,
			"/register",
			bytes.NewBufferString(`{"socker_addr": "http://localhost:1234"}`),
		)
		c.Request = req
		ctrl.Register(c)
		assert.Equal(t, http.StatusOK, w.Code)

	})
}

func TestSyncController_Unregister(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	sync := mock.NewMockSyncService(mockCtrl)
	ctrl := NewSyncController(sync)
	assert.NotNil(t, ctrl)

	t.Run("no_username", func(t *testing.T) {
		// Test case 1: Missing username
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.POST("/unregister", ctrl.Unregister)
		req, _ := http.NewRequest(http.MethodPost, "/unregister", nil)
		c.Request = req
		ctrl.Unregister(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
	t.Run("invalid_data", func(t *testing.T) {
		// Test case 2: Invalid JSON payload
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.POST("/unregister", ctrl.Unregister)
		c.Set(middleware.UsernameContextValue, "blokhinnv")
		req, _ := http.NewRequest(
			http.MethodPost,
			"/unregister",
			bytes.NewBufferString(`{"invalid": json}`),
		)
		c.Request = req
		ctrl.Unregister(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ok", func(t *testing.T) {
		// Test case 3: Valid registration
		sync.EXPECT().Unregister(gomock.Any()).Times(1)
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.POST("/unregister", ctrl.Unregister)
		c.Set(middleware.UsernameContextValue, "blokhinnv")
		req, _ := http.NewRequest(
			http.MethodPost,
			"/unregister",
			bytes.NewBufferString(`{"socker_addr": "http://localhost:1234"}`),
		)
		c.Request = req
		ctrl.Unregister(c)
		assert.Equal(t, http.StatusOK, w.Code)

	})
}

func TestNewSyncController(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	sync := mock.NewMockSyncService(mockCtrl)
	ctrl := NewSyncController(sync)
	assert.NotNil(t, ctrl)
}
