package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/blokhinnv/gophkeeper/internal/server/errors"
	"github.com/blokhinnv/gophkeeper/internal/server/service/mock"
)

func TestNewAuthController(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// create a new authController instance with a mocked service
	srvc := mock.NewMockAuthService(mockCtrl)
	ctrl := NewAuthController(srvc)
	assert.NotNil(t, ctrl)
}

func TestAuthController_Register(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// create a new authController instance with a mocked service
	srvc := mock.NewMockAuthService(mockCtrl)
	ctrl := NewAuthController(srvc)
	// create a valid user credentials JSON
	userJSON := `{"username": "testuser", "password": "testpassword"}`
	t.Run("ok", func(t *testing.T) {
		srvc.EXPECT().
			Register(gomock.Eq("testuser"), gomock.Eq("testpassword")).
			Times(1).
			Return(nil)

		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		// test registering a new user
		r.POST("/register", ctrl.Register)
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(userJSON))
		c.Request = req
		ctrl.Register(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("duplicate", func(t *testing.T) {
		// test registering a duplicate username
		srvc.EXPECT().
			Register(gomock.Eq("testuser"), gomock.Eq("testpassword")).
			Times(1).
			Return(errors.ErrUsernameIsTakenMongo)
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.POST("/register", ctrl.Register)
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(userJSON))
		c.Request = req
		ctrl.Register(c)
		assert.Equal(t, http.StatusConflict, w.Code)
	})
	t.Run("other_db_error", func(t *testing.T) {
		// test registering a duplicate username
		srvc.EXPECT().
			Register(gomock.Eq("testuser"), gomock.Eq("testpassword")).
			Times(1).
			Return(fmt.Errorf("some db error"))
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.POST("/register", ctrl.Register)
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(userJSON))
		c.Request = req
		ctrl.Register(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("invalid_data", func(t *testing.T) {
		// test invalid JSON data
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.POST("/register", ctrl.Register)
		req, _ := http.NewRequest(
			http.MethodPost,
			"/register",
			bytes.NewBufferString(`{"foo": "bar"}`),
		)
		c.Request = req
		ctrl.Register(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAuthController_Login(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// create a new authController instance with a mocked service
	srvc := mock.NewMockAuthService(mockCtrl)
	ctrl := NewAuthController(srvc)
	// create a valid user credentials JSON
	t.Run("ok", func(t *testing.T) {
		// test logging in with valid credentials
		srvc.EXPECT().
			Login(gomock.Eq("testuser"), gomock.Eq("testpassword")).
			Times(1).
			Return("some-token", nil)

		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.POST("/login", ctrl.Login)
		req, _ := http.NewRequest(
			http.MethodPost,
			"/login",
			bytes.NewBufferString(`{"username": "testuser", "password": "testpassword"}`),
		)
		c.Request = req
		ctrl.Login(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("duplicate", func(t *testing.T) {
		// test logging in with invalid credentials
		srvc.EXPECT().
			Login(gomock.Eq("testuser"), gomock.Eq("wrongpassword")).
			Times(1).
			Return("", errors.ErrNoDocuments)
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.POST("/login", ctrl.Login)
		req, _ := http.NewRequest(
			http.MethodPost,
			"/login",
			bytes.NewBufferString(`{"username": "testuser", "password": "wrongpassword"}`),
		)
		c.Request = req
		ctrl.Login(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
	t.Run("invalid_data", func(t *testing.T) {
		// test invalid JSON data
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.POST("/login", ctrl.Login)
		req, _ := http.NewRequest(
			http.MethodPost,
			"/login",
			bytes.NewBufferString(`{"foo":}`),
		)
		c.Request = req
		ctrl.Login(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
