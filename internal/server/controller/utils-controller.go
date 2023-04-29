package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/blokhinnv/gophkeeper/internal/server/service"
)

// UtilsController defines the interface for utility controller,
// which includes methods for testing server availability.
type UtilsController interface {
	// Ping tests server availability.
	Ping(*gin.Context)
}

// utilsController implements UtilsController interface.
type utilsController struct {
	service service.UtilsService
}

// NewUtilsController creates a new instance of UtilsController.
func NewUtilsController(service service.UtilsService) UtilsController {
	return &utilsController{
		service: service,
	}
}

// Ping godoc
//
//	@Summary Ping server
//	@Description Returns plain text response with a "pong" message if the server is available, otherwise returns an error message.
//	@Produce plain
//	@ID Ping
//	@Tags Utils
//	@Success 200 {string}	string	"pong"
//	@Failure 500 {string}	string	"not pong"
//	@Router /api/ping [get]
func (c *utilsController) Ping(ctx *gin.Context) {
	err := c.service.Ping()
	if err != nil {
		ctx.String(http.StatusInternalServerError, "not pong")
		return
	}
	ctx.String(http.StatusOK, "pong")
}
