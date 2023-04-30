package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/blokhinnv/gophkeeper/internal/server/errors"
	"github.com/blokhinnv/gophkeeper/internal/server/models"
	"github.com/blokhinnv/gophkeeper/pkg/encrypt"
)

// StorageService is an interface that defines the methods to store and retrieve untyped records.
type StorageService interface {
	// Store stores a new untyped record in a specified collection.
	Store(
		ctx context.Context,
		collectionName models.Collection,
		record models.UntypedRecord,
	) (string, error)
	// GetAll retrieves all untyped records for a specified collection and username.
	GetAll(
		ctx context.Context,
		collectionName models.Collection,
		username string,
	) ([]models.UntypedRecord, error)
	// Updates the data and metadata of the document.
	Update(
		ctx context.Context,
		collectionName models.Collection,
		username string,
		id primitive.ObjectID,
		newData any,
		newMetadata models.Metadata,
	) error
	// Deletes the document from collection.
	Delete(
		ctx context.Context,
		collectionName models.Collection,
		username string,
		id primitive.ObjectID,
	) error
}

// storageService is a struct that implements the StorageService
// interface and uses MongoDB for data storage.
type storageService struct {
	db            *mongo.Database
	encryptionKey string
}

// NewStorageService creates a new storageService instance.
func NewStorageService(db *mongo.Database, encryptionKey string) StorageService {
	return &storageService{
		db:            db,
		encryptionKey: encryptionKey,
	}
}

// Store stores a new untyped record in a specified collection.
func (t *storageService) Store(
	ctx context.Context,
	collectionName models.Collection,
	record models.UntypedRecord,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	collection := t.db.Collection(string(collectionName))

	var encryptedData any
	var err error
	switch collectionName {
	case models.TextCollection:
		encryptedData, err = encrypt.EncryptString(record.Data.(string), t.encryptionKey)
	default:
		encryptedData, err = encrypt.EncryptMap(record.Data.(map[string]any), t.encryptionKey)
	}
	if err != nil {
		return "", err
	}

	res, err := collection.InsertOne(ctx, bson.D{
		{Key: "username", Value: record.Username},
		{Key: "data", Value: encryptedData},
		{Key: "metadata", Value: record.Metadata},
	})
	stringObjectID := res.InsertedID.(primitive.ObjectID).Hex()
	return stringObjectID, err
}

// GetAll retrieves all untyped records for a specified collection and username.
func (t *storageService) GetAll(
	ctx context.Context,
	collectionName models.Collection,
	username string,
) ([]models.UntypedRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	result := make([]models.UntypedRecord, 0)
	collection := t.db.Collection(string(collectionName))
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
		if err != nil {
			return nil, err
		}
		var decryptedData any
		v, ok := r.Data.(bson.D)
		if ok {
			decryptedData, err = encrypt.DecryptMap(v.Map(), t.encryptionKey)
		} else {
			decryptedData, err = encrypt.DecryptString(r.Data.(string), t.encryptionKey)
		}
		r.Data = decryptedData
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
func (t *storageService) Update(
	ctx context.Context,
	collectionName models.Collection,
	username string,
	id primitive.ObjectID,
	newData any,
	newMetadata models.Metadata,
) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var encryptedNewData any
	var err error
	switch collectionName {
	case models.TextCollection:
		encryptedNewData, err = encrypt.EncryptString(newData.(string), t.encryptionKey)
	default:
		encryptedNewData, err = encrypt.EncryptMap(newData.(map[string]any), t.encryptionKey)
	}
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id, "username": username}
	upd := bson.D{{
		Key: "$set",
		Value: bson.D{
			{Key: "data", Value: encryptedNewData},
			{Key: "metadata", Value: newMetadata},
		},
	}}
	collection := t.db.Collection(string(collectionName))
	res, err := collection.UpdateOne(ctx, filter, upd)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.ErrRecordNotFound
	}
	return nil
}

// Delete removes a document from the specified collection using its ObjectID.
func (t *storageService) Delete(
	ctx context.Context,
	collectionName models.Collection,
	username string,
	id primitive.ObjectID,
) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	filter := bson.M{"_id": id, "username": username}
	collection := t.db.Collection(string(collectionName))
	res, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.ErrRecordNotFound
	}
	return nil
}
