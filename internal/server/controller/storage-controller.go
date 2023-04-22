package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"

	"gophkeeper/internal/server/errors"
	"gophkeeper/internal/server/middleware"
	"gophkeeper/internal/server/models"
	"gophkeeper/internal/server/service"
	"gophkeeper/internal/server/validation"
)

// StorageController defines the interface for storage
// controller, which includes methods for storing data and getting all data.
type StorageController interface {
	// Store saves an untyped record to the database based on the data provided in the request.
	Store(ctx *gin.Context)
	// GetAll returns all the untyped records from the database based on the data provided in the request.
	GetAll(ctx *gin.Context)
}

// storageController implements StorageController interface.
type storageController struct {
	service service.StorageService
}

// NewStorageController creates a new instance of StorageController with the given StorageService.
func NewStorageController(service service.StorageService) StorageController {
	return &storageController{
		service: service,
	}
}

// getCollectionName returns the collection name based on the URL of the request.
func (c *storageController) getCollectionName(url string) string {
	switch {
	case strings.Contains(url, "text"):
		return models.TextCollection
	case strings.Contains(url, "binary"):
		return models.BinaryCollection
	case strings.Contains(url, "cards"):
		return models.CardCollection
	case strings.Contains(url, "credentials"):
		return models.CredentialsCollection
	default:
		return models.TextCollection
	}
}

// validateDataField validates the data field of the untyped record based on the collection name.
func (c *storageController) validateDataField(data any, collectionName string) error {
	// TODO: check binary

	switch collectionName {
	case models.CredentialsCollection:
		var c models.Credential
		err := mapstructure.Decode(data, &c)
		if err != nil {
			return err
		}

		err = validation.Validate.Struct(c)
		return err
	case models.CardCollection:
		var c models.CardInfo
		err := mapstructure.Decode(data, &c)
		if err != nil {
			return err
		}
		err = validation.Validate.Struct(c)
		return err
	default:
		return nil
	}

}

// Store saves an untyped record to the database based on the data provided in the request.
func (c *storageController) Store(ctx *gin.Context) {
	username := ctx.GetString(middleware.UsernameContextValue)
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrNoUsernameProvided.Error()})
		return
	}
	record := models.UntypedRecord{
		Username: username,
	}
	if err := ctx.ShouldBindJSON(&record); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	collectionName := c.getCollectionName(ctx.Request.RequestURI)
	if err := c.validateDataField(record.Data, collectionName); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.service.Store(collectionName, record)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, nil)
}

// GetAll returns all the untyped records from the database based on the data provided in the request.
func (c *storageController) GetAll(ctx *gin.Context) {
	username := ctx.GetString(middleware.UsernameContextValue)
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrNoUsernameProvided.Error()})
		return
	}
	collectionName := c.getCollectionName(ctx.Request.RequestURI)
	records, err := c.service.GetAll(collectionName, username)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, records)
}
