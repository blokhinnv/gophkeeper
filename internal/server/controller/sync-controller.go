package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gophkeeper/internal/server/errors"
	"gophkeeper/internal/server/middleware"
	"gophkeeper/internal/server/models"
	"gophkeeper/internal/server/service"
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

// Register registers a new client with the server.
func (s *syncController) Register(ctx *gin.Context) {
	username := ctx.GetString(middleware.UsernameContextValue)
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrNoUsernameProvided.Error()})
		return
	}
	client := &models.Client{
		Addr:     ctx.Request.RemoteAddr,
		Username: username,
	}
	s.service.Register(client)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "registered",
		"info": gin.H{
			"addr":     client.Addr,
			"username": client.Username,
		},
	})
}

// Unregister unregisters an existing client from the server.
func (s *syncController) Unregister(ctx *gin.Context) {
	username := ctx.GetString(middleware.UsernameContextValue)
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrNoUsernameProvided.Error()})
		return
	}
	client := &models.Client{
		Addr:     ctx.Request.RemoteAddr,
		Username: username,
	}
	s.service.Unregister(client)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "unregistered",
		"info": gin.H{
			"addr":     client.Addr,
			"username": client.Username,
		},
	})
}
