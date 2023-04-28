// Package controller provides implementations for various
// controllers used by the server to handle client
// requests related to authentication,
// authorization, and other functionalities.
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/blokhinnv/gophkeeper/internal/server/errors"
	"github.com/blokhinnv/gophkeeper/internal/server/models"
	"github.com/blokhinnv/gophkeeper/internal/server/service"
)

// AuthController defines the interface for authentication controller.
type AuthController interface {
	// Register saves a new user to the database based on the data provided in the request.
	Register(*gin.Context)
	// Login checks users' credentials and returns a JWT token.
	Login(*gin.Context)
}

// authController implements AuthController interface.
type authController struct {
	service service.AuthService
}

// NewAuthController creates a new instance of AuthController.
func NewAuthController(service service.AuthService) AuthController {
	return &authController{
		service: service,
	}
}

// Register is the controller method for registering a new user.
func (c *authController) Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.service.Register(user.Username, user.Password); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			ctx.JSON(http.StatusConflict, gin.H{"error": errors.ErrUsernameIsTaken.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "registrated"})
}

// Login is the controller method for user login.
func (c *authController) Login(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tok, err := c.service.Login(user.Username, user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "logged in", "tok": tok})
}
