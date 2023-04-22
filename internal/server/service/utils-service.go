package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// UtilsService is an interface that defines a method for pinging a MongoDB database.
type UtilsService interface {
	Ping() error
}

// utilsService is a concrete implementation of the UtilsService interface.
type utilsService struct {
	client *mongo.Client
}

// NewUtilsService returns a new instance of the UtilsService interface.
// It takes a pointer to a mongo.Client instance as its only parameter.
func NewUtilsService(
	client *mongo.Client,
) UtilsService {
	return &utilsService{
		client: client,
	}
}

// Ping pings the MongoDB database and returns an error if the ping fails.
func (t *utilsService) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return t.client.Ping(ctx, readpref.Primary())
}
