package controller

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/blokhinnv/gophkeeper/internal/server/service/mock"
)

func TestUtilsController_Ping(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	utils := mock.NewMockUtilsService(mockCtrl)
	ctrl := NewUtilsController(utils)
	assert.NotNil(t, ctrl)

	t.Run("ok", func(t *testing.T) {
		// Test case 1: Successful ping
		utils.EXPECT().Ping().Times(1).Return(nil)
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.GET("/ping", ctrl.Ping)
		req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
		c.Request = req
		ctrl.Ping(c)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "pong", w.Body.String())
	})

	t.Run("not_ok", func(t *testing.T) {
		// Test case 2: Failed ping
		utils.EXPECT().Ping().Times(1).Return(fmt.Errorf("some error"))
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.GET("/ping", ctrl.Ping)
		req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
		c.Request = req
		ctrl.Ping(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "not pong", w.Body.String())
	})
}

func TestNewUtilsController(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	utils := mock.NewMockUtilsService(mockCtrl)
	ctrl := NewUtilsController(utils)
	assert.NotNil(t, ctrl)
}
