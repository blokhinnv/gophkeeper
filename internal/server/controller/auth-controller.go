// Package controller provides implementations for various
// controllers used by the server to handle client
// requests related to authentication,
// authorization, and other functionalities.
package controller

import (
	"fmt"
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

// Register godoc
//
//	@Summary Register a new user
//	@Description Register a new user with provided credentials
//	@Produce plain
//	@ID Register
//	@Tags Authy
//	@Param	credentials body	models.UserCredentials	true	"Credentials"
//	@Success 200 {string}	string	"success"
//	@Failure 400 {string}	string	"Bad Request"
//	@Failure 409 {string}	string	"Username is already taken"
//	@Router /api/user/register [put]
func (c *authController) Register(ctx *gin.Context) {
	var user models.UserCredentials
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	if err := c.service.Register(user.Username, user.Password); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			ctx.String(http.StatusConflict,
				fmt.Sprintf(
					"%v: %v",
					errors.ErrUsernameIsTaken.Error(),
					user.Username,
				),
			)
			return
		}
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.String(http.StatusOK, "success")
}

// Login godoc
//
//	@Summary Logs in a user
//	@Description Logs in a user with the provided username and password
//	@Produce plain
//	@ID Login
//	@Tags Authy
//	@Param	credentials body	models.UserCredentials	true	"Credentials"
//	@Success 200 {string}	string	"some JWT token"
//	@Failure 400 {string}	string	"no username provided"
//	@Failure 401 {string}	string	"username or password is incorrect: testuser/qwerty"
//	@Router /api/user/login [put]
func (c *authController) Login(ctx *gin.Context) {
	var user models.UserCredentials
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	tok, err := c.service.Login(user.Username, user.Password)
	if err != nil {
		ctx.String(
			http.StatusUnauthorized,
			fmt.Sprintf(
				"%v: %v %v",
				errors.ErrBadCredentials.Error(),
				user.Username,
				user.Password,
			),
		)
		return
	}
	ctx.String(http.StatusOK, tok)
}
