package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/blokhinnv/gophkeeper/internal/server/auth"
)

func TestJWTAuthMiddleware(t *testing.T) {
	signingKey := []byte("secret")
	t.Run("unauthorized", func(t *testing.T) {
		// Test unauthorized request
		// Create a mock gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest("GET", "/test", nil)
		c.Request = req
		JWTAuthMiddleware(signingKey)(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
	t.Run("authorized", func(t *testing.T) {
		// Test authorized request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest("GET", "/test", nil)
		c.Request = req
		tokenString, _ := auth.GenerateJWTToken("user", signingKey, time.Hour)
		req.Header.Set("Authorization", "Bearer: "+tokenString)
		JWTAuthMiddleware(signingKey)(c)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "user", c.GetString(UsernameContextValue))
	})
}
