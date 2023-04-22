// Package middleware contains different custom middlewares.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gophkeeper/internal/server/auth"
)

// UsernameContextValue is the key used to set and get the username value in gin.Context.
const UsernameContextValue = "username"

// JWTAuthMiddleware is a middleware that performs JWT token validation.
// It returns a gin.HandlerFunc which can be used in a gin route.
func JWTAuthMiddleware(signingKey []byte) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")
		var tokenString string
		if len(strings.Split(authHeader, " ")) == 2 {
			tokenString = strings.Split(authHeader, " ")[1]
		}

		username, err := auth.ValidateJWTToken(tokenString, signingKey)
		if err != nil {
			ctx.String(http.StatusUnauthorized, "Unauthorized")
			ctx.Abort()
			return
		}
		ctx.Set(UsernameContextValue, username)
		ctx.Next()
	}
}
