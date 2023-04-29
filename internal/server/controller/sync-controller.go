package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/blokhinnv/gophkeeper/internal/server/errors"
	"github.com/blokhinnv/gophkeeper/internal/server/middleware"
	"github.com/blokhinnv/gophkeeper/internal/server/models"
	"github.com/blokhinnv/gophkeeper/internal/server/service"
)

// SyncController defines the interface for handling sync-related HTTP requests.
type SyncController interface {
	Register(ctx *gin.Context)
	Unregister(ctx *gin.Context)
}

// syncController implements the SyncController interface.
type syncController struct {
	service service.SyncService
}

// NewSyncController creates a new SyncController instance with the given SyncService.
func NewSyncController(service service.SyncService) SyncController {
	return &syncController{
		service: service,
	}
}

// Register godoc
//
//	@Summary Registers a new client with the synchronization service.
//	@Security bearerAuth
//	@Description Allows a client to register with the synchronization service.
//	@Accept json
//	@Produce plain
//	@ID RegisterClient
//	@Tags Sync
//	@Param	client	body	models.Client	true	"Client"
//	@Success 200 {string}	string	"client registered"
//	@Failure 400 {string}	string	"Bad Request"
//	@Failure 401 {string}	string	"No username provided"
//	@Router /api/sync/register [post]
func (s *syncController) Register(ctx *gin.Context) {
	username := ctx.GetString(middleware.UsernameContextValue)
	if username == "" {
		ctx.String(http.StatusUnauthorized, errors.ErrNoUsernameProvided.Error())
		return
	}
	client := &models.Client{
		Username: username,
	}
	if err := ctx.ShouldBindJSON(&client); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	s.service.Register(client)
	ctx.String(http.StatusOK, "client registered")
}

// Unregister
// Store godoc
//
//	@Summary Unregisters an existing client from the server.
//	@Security bearerAuth
//	@Description Allows a client to unregister from the synchronization service.
//	@Accept json
//	@Produce plain
//	@ID UnregisterClient
//	@Tags Sync
//	@Param	client	body	models.Client	true	"Client"
//	@Success 200 {string}	string	"client unregistered"
//	@Failure 400 {string}	string	"Bad Request"
//	@Failure 401 {string}	string	"No username provided"
//	@Router /api/sync/unregister [post]
func (s *syncController) Unregister(ctx *gin.Context) {
	username := ctx.GetString(middleware.UsernameContextValue)
	if username == "" {
		ctx.String(http.StatusBadRequest, errors.ErrNoUsernameProvided.Error())
		return
	}
	client := &models.Client{
		Username: username,
	}
	if err := ctx.ShouldBindJSON(&client); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	s.service.Unregister(client)
	ctx.String(http.StatusOK, "client unregistered")
}
