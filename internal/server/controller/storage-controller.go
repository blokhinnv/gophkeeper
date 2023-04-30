package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson/primitive"

	srvErrors "github.com/blokhinnv/gophkeeper/internal/server/errors"
	"github.com/blokhinnv/gophkeeper/internal/server/middleware"
	"github.com/blokhinnv/gophkeeper/internal/server/models"
	"github.com/blokhinnv/gophkeeper/internal/server/service"
	"github.com/blokhinnv/gophkeeper/internal/server/validation"
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
	sync    service.SyncService
}

// NewStorageController creates a new instance of StorageController with the given StorageService.
func NewStorageController(
	service service.StorageService,
	sync service.SyncService,
) StorageController {
	return &storageController{
		service: service,
		sync:    sync,
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
		var c models.BinaryInfo
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

// Store godoc
//
//	@Summary Store an untyped record to the database.
//	@Security bearerAuth
//	@Description Stores an untyped record to the database based on the data provided in the request.
//	@Accept json
//	@Produce plain
//	@ID Store
//	@Tags Storage
//	@Param	record	body	models.UntypedRecordContent	true	"Record"
//	@Param        collectionName   path      string  true  "Collection name"
//	@Success 202 {string}	string	"Record added to collection"
//	@Failure 400 {string}	string	"Bad Request"
//	@Failure 401 {string}	string	"No username provided"
//	@Router /api/store/{collectionName} [put]
func (c *storageController) Store(ctx *gin.Context) {
	username := ctx.GetString(middleware.UsernameContextValue)
	if username == "" {
		ctx.String(http.StatusUnauthorized, srvErrors.ErrNoUsernameProvided.Error())
		return
	}
	record := models.UntypedRecord{
		Username: username,
	}
	if err := ctx.ShouldBindJSON(&record.UntypedRecordContent); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	collectionName, err := models.NewCollection(ctx.Param("collectionName"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	if err := c.validateDataField(record.Data, collectionName); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	id, err := c.service.Store(ctx.Request.Context(), collectionName, record)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	go c.sync.Signal(&models.Client{Username: username})
	ctx.String(
		http.StatusAccepted,
		fmt.Sprintf(
			"Record added to %v collection: id=%v data=%v metadata=%v",
			collectionName,
			id,
			record.Data,
			record.Metadata,
		),
	)
}

// GetAll godoc
//
//	@Summary Retrieve all untyped records for the authenticated user from a collection.
//	@Description Returns all the untyped records from the database based on the data provided in the request.
//	@Security bearerAuth
//	@Accept json
//	@Produce json
//	@ID GetAll
//	@Tags Storage
//	@Param        collectionName   path      string  true  "Collection name"
//	@Success 200 {array}	models.UntypedRecord	"Record added by the user in the specified collection"
//	@Failure 400 {string}	string	"Bad Request"
//	@Failure 401 {string}	string	"No username provided"
//	@Router /api/store/{collectionName} [get]
func (c *storageController) GetAll(ctx *gin.Context) {
	username := ctx.GetString(middleware.UsernameContextValue)
	if username == "" {
		ctx.String(http.StatusUnauthorized, srvErrors.ErrNoUsernameProvided.Error())
		return
	}
	collectionName, err := models.NewCollection(ctx.Param("collectionName"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	records, err := c.service.GetAll(ctx.Request.Context(), collectionName, username)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, records)
}

// Update godoc
//
//	@Summary Update an existing record in the database.
//	@Description Updates the data and metadata of a document in the collection specified by the request URL, based on the data provided in the request body. The updated document is identified by its ID, which is included in the request body as well.
//	@Security bearerAuth
//	@Accept json
//	@Produce plain
//	@ID Update
//	@Tags Storage
//	@Param	record	body	models.UntypedRecord	true	"Record"
//	@Param        collectionName   path      string  true  "Collection name"
//	@Success 202 {string}	string	"Record updated"
//	@Failure 400 {string}	string	"Bad Request"
//	@Failure 401 {string}	string	"No username provided"
//	@Failure 500 {string}	string	"Server error"
//	@Router /api/store/{collectionName} [post]
func (c *storageController) Update(ctx *gin.Context) {
	username := ctx.GetString(middleware.UsernameContextValue)
	if username == "" {
		ctx.String(http.StatusUnauthorized, srvErrors.ErrNoUsernameProvided.Error())
		return
	}
	record := models.UntypedRecord{
		Username: username,
	}
	if err := ctx.ShouldBindJSON(&record); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	collectionName, err := models.NewCollection(ctx.Param("collectionName"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	if err := c.validateDataField(record.Data, collectionName); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.Update(
		ctx.Request.Context(),
		collectionName,
		username,
		record.RecordID,
		record.Data,
		record.Metadata,
	)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, srvErrors.ErrRecordNotFound) {
			status = http.StatusBadRequest
		}
		ctx.String(status, err.Error())
		return
	}
	go c.sync.Signal(&models.Client{Username: username})
	ctx.String(
		http.StatusAccepted,
		fmt.Sprintf(
			"Record id=%v updated in %v collection: data=%v metadata=%v",
			record.RecordID,
			collectionName,
			record.Data,
			record.Metadata,
		),
	)
}

type deleteRequestBody struct {
	RecordID primitive.ObjectID `json:"record_id" bson:"_id" binding:"required"`
}

// Delete godoc
//
//	@Summary Delete a record by ID
//	@Description Deletes a record from the specified collection by ID.
//	@Security bearerAuth
//	@Accept json
//	@Produce plain
//	@ID Delete
//	@Tags Storage
//	@Param	record_id	body	deleteRequestBody	true	"RecordID"
//	@Param        collectionName   path      string  true  "Collection name"
//	@Success 200 {string}	string	"Record deleted"
//	@Failure 400 {string}	string	"Bad Request"
//	@Failure 401 {string}	string	"No username provided"
//	@Failure 500 {string}	string	"Server error"
//	@Router /api/store/{collectionName} [delete]
func (c *storageController) Delete(ctx *gin.Context) {
	username := ctx.GetString(middleware.UsernameContextValue)
	if username == "" {
		ctx.String(http.StatusUnauthorized, srvErrors.ErrNoUsernameProvided.Error())
		return
	}
	record := deleteRequestBody{}

	if err := ctx.ShouldBindJSON(&record); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	collectionName, err := models.NewCollection(ctx.Param("collectionName"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	err = c.service.Delete(ctx.Request.Context(), collectionName, username, record.RecordID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, srvErrors.ErrRecordNotFound) {
			status = http.StatusBadRequest
		}
		ctx.String(status, err.Error())
		return
	}
	go c.sync.Signal(&models.Client{Username: username})
	ctx.String(
		http.StatusOK,
		fmt.Sprintf("Record id=%v deleted from %v collection", record.RecordID, collectionName),
	)
}
