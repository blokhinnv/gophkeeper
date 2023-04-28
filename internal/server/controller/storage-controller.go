package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson/primitive"

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
	// Updates the data and metadata of a document in the collection specified by the request URL.
	Update(ctx *gin.Context)
	// Delete deletes a record from the collection specified in the request URI.
	Delete(ctx *gin.Context)
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
func (c *storageController) getCollectionName(url string) models.Collection {
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
func (c *storageController) validateDataField(
	data any,
	collectionName models.Collection,
) error {
	switch collectionName {
	case models.CredentialsCollection:
		var c models.CredentialInfo
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
	case models.BinaryCollection:
		return validation.Validate.Var(data.(string), "base64")
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

	err := c.service.Store(ctx.Request.Context(), collectionName, record)
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
	records, err := c.service.GetAll(ctx.Request.Context(), collectionName, username)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, records)
}

// Updates the data and metadata of a document in the collection specified by the request URL,
// based on the data provided in the request body. The updated document is identified by its ID,
// which is included in the request body as well.
func (c *storageController) Update(ctx *gin.Context) {
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

	err := c.service.Update(
		ctx.Request.Context(),
		collectionName,
		username,
		record.RecordID,
		record.Data,
		record.Metadata,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// Delete deletes a record from the collection specified in the request URI
// using the provided record ID and the authenticated user's username.
func (c *storageController) Delete(ctx *gin.Context) {
	username := ctx.GetString(middleware.UsernameContextValue)
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrNoUsernameProvided.Error()})
		return
	}
	record := struct {
		RecordID primitive.ObjectID `json:"record_id" bson:"_id" binding:"required"`
	}{}

	if err := ctx.ShouldBindJSON(&record); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	collectionName := c.getCollectionName(ctx.Request.RequestURI)

	err := c.service.Delete(ctx.Request.Context(), collectionName, username, record.RecordID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
