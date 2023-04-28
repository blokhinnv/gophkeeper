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

// Ping returns a JSON response with a "pong" message
// if the server is available, otherwise returns an error message.
func (c *utilsController) Ping(ctx *gin.Context) {
	err := c.service.Ping()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "not pong",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
