package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConnection struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoDBConnection(ctx context.Context, uri, dbName string) (*MongoDBConnection, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &MongoDBConnection{
		client:   client,
		database: client.Database(dbName),
	}, nil
}

func (mc *MongoDBConnection) Close(ctx context.Context) error {
	return mc.client.Disconnect(ctx)
}

func (mc *MongoDBConnection) Database() *mongo.Database {
	return mc.database
}
