package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"gophkeeper/internal/server/models"
)

// StorageService is an interface that defines the methods to store and retrieve untyped records.
type StorageService interface {
	// Store stores a new untyped record in a specified collection.
	Store(collectionName string, record models.UntypedRecord) error
	// GetAll retrieves all untyped records for a specified collection and username.
	GetAll(collectionName, username string) ([]models.UntypedRecord, error)
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
//   - collectionName: The name of the collection to store the record in.
//   - record: The untyped record to store.
//
// Returns:
//   - error: An error if the store operation failed, or nil if successful.
func (t *storageService) Store(collectionName string, record models.UntypedRecord) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
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
//   - collectionName: The name of the collection to retrieve records from.
//   - username: The login of a user which retrieves the records.

func (t *storageService) GetAll(
	collectionName, username string,
) ([]models.UntypedRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
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
