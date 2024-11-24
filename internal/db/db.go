package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// New creates and configures a new connection to MongoDB
func New(addr string, maxConn int, maxIdleTime string) (*mongo.Client, error) {
	// Create a new client option for MongoDB connection
	clientOptions := options.Client().ApplyURI(addr)

	// Set max pool size and idle time
	clientOptions.SetMaxPoolSize(uint64(maxConn))

	// Parse maxIdleTime into a duration
	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	clientOptions.SetMaxConnIdleTime(duration)

	// Create a context with a timeout for establishing the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Verify the connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("could not connect to MongoDB: %v", err)
	}

	return client, nil
}
