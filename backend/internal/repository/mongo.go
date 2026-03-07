package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const connectTimeout = 10 * time.Second

type DB struct {
	client *mongo.Client
	name   string
}

func NewMongoDB(uri, dbName string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo connect: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongo ping: %w", err)
	}

	return &DB{client: client, name: dbName}, nil
}

func (db *DB) Collection(name string) *mongo.Collection {
	return db.client.Database(db.name).Collection(name)
}

func (db *DB) Close(ctx context.Context) error {
	return db.client.Disconnect(ctx)
}
