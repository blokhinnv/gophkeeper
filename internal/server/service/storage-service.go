package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"gophkeeper/internal/server/models"
)

// StorageService is an interface that defines the methods to store and retrieve untyped records.
type StorageService interface {
	// Store stores a new untyped record in a specified collection.
	Store(ctx context.Context, collectionName string, record models.UntypedRecord) error
	// GetAll retrieves all untyped records for a specified collection and username.
	GetAll(ctx context.Context, collectionName, username string) ([]models.UntypedRecord, error)
	// Updates the data and metadata of the document.
	Update(
		ctx context.Context,
		collectionName, username string,
		id primitive.ObjectID,
		newData any,
		newMetadata models.Metadata,
	) error
	// TODO:
	Delete(ctx context.Context, collectionName, username string, id primitive.ObjectID) error
}

// storageService is a struct that implements the StorageService
// interface and uses MongoDB for data storage.
type storageService struct {
	db *mongo.Database
}

// NewStorageService creates a new storageService instance.
// Parameters:
//   - db: A pointer to a mongo.Database instance.
//
// Returns:
//   - StorageService: A new instance of the storageService struct.
func NewStorageService(db *mongo.Database) StorageService {
	return &storageService{
		db: db,
	}
}

// Store stores a new untyped record in a specified collection.
// Parameters:
// - ctx (context.Context): The context to use for the operation.
// - collectionName: The name of the collection to store the record in.
// - record: The untyped record to store.
//
// Returns:
//   - error: An error if the store operation failed, or nil if successful.
func (t *storageService) Store(
	ctx context.Context,
	collectionName string,
	record models.UntypedRecord,
) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	collection := t.db.Collection(collectionName)
	_, err := collection.InsertOne(ctx, bson.D{
		{Key: "username", Value: record.Username},
		{Key: "data", Value: record.Data},
		{Key: "metadata", Value: record.Metadata},
	})
	return err
}

// GetAll retrieves all untyped records for a specified collection and username.
// Parameters:
// - ctx (context.Context): The context to use for the operation.
// - collectionName: The name of the collection to retrieve records from.
// - username: The login of a user which retrieves the records.
func (t *storageService) GetAll(
	ctx context.Context,
	collectionName, username string,
) ([]models.UntypedRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	result := make([]models.UntypedRecord, 0)
	collection := t.db.Collection(collectionName)
	cur, err := collection.Find(
		ctx,
		bson.D{
			{Key: "username", Value: username},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var r models.UntypedRecord
		err := cur.Decode(&r)
		fmt.Println(r)
		v, ok := r.Data.(bson.D)
		if ok {
			r.Data = v.Map()
		}
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// Updates the data and metadata of the document with the specified ID in the
// collection with the specified name, using the new data and metadata values.
//
// Parameters:
// - ctx (context.Context): The context to use for the operation.
// - collectionName: string, the name of the collection to update the document in.
// - username: string, the username of the user performing the update operation.
// - id: primitive.ObjectID, the ID of the document to update.
// - newData: any, the new data value to set for the "data" field of the document.
// - newMetadata: models.Metadata, the new metadata value to set for the "metadata" field of the document.
//
// Returns:
//   - error: an error, if any occurred during the update operation (e.g., if the document
//     could not be found or the update failed).
func (t *storageService) Update(
	ctx context.Context,
	collectionName, username string,
	id primitive.ObjectID,
	newData any,
	newMetadata models.Metadata,
) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	filter := bson.M{"_id": id}
	upd := bson.D{{
		Key: "$set",
		Value: bson.D{
			{Key: "data", Value: newData},
			{Key: "metadata", Value: newMetadata},
		},
	}}
	collection := t.db.Collection(collectionName)
	_, err := collection.UpdateOne(ctx, filter, upd)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a document from the specified collection using its ObjectID.
//
// Parameters:
// - ctx (context.Context): The context to use for the operation.
// - collectionName (string): The name of the collection to delete the document from.
// - username (string): The username of the user attempting to delete the document.
// - id (primitive.ObjectID): The ObjectID of the document to delete.
//
// Returns:
// - error: Any error that occurred during the operation, or nil if successful.
func (t *storageService) Delete(
	ctx context.Context,
	collectionName, username string,
	id primitive.ObjectID,
) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	filter := bson.M{"_id": id}
	collection := t.db.Collection(collectionName)
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
