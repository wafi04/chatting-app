package database

import (
	"context"
	"fmt"
	"time"

	"github.com/wafi04/chatting-app/config/env"
	"github.com/wafi04/chatting-app/services/shared/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDB represents the MongoDB connection handler
type MongoDB struct {
	Client *mongo.Client
	log    logger.Logger
}

// ConnectMongoDB establishes a connection to MongoDB
func ConnectMongoDB(log *logger.Logger) (*MongoDB, error) {
	// Load MongoDB URI from environment variables
	uri := env.LoadEnv("MONGO_URL")
	if uri == "" {
		return nil, fmt.Errorf("MONGO_URL environment variable is not set")
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configure client options
	clientOptions := options.Client().
		ApplyURI(uri).
		SetConnectTimeout(5 * time.Second).
		SetServerSelectionTimeout(5 * time.Second).
		SetSocketTimeout(30 * time.Second)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Verify the connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// Log successful connection
	log.Log(logger.InfoLevel, "Connected to MongoDB")

	// Return the MongoDB instance
	return &MongoDB{
		Client: client,
	}, nil
}

// Close closes the MongoDB connection
func (d *MongoDB) Close() error {
	if d.Client == nil {
		d.log.Log(logger.WarnLevel, "MongoDB client is already closed or not initialized")
		return nil
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Disconnect from MongoDB
	err := d.Client.Disconnect(ctx)
	if err != nil {
		d.log.Log(logger.ErrorLevel, "failed to disconnect from MongoDB: %v", err)
		return err
	}

	// Log successful disconnection
	d.log.Log(logger.InfoLevel, "Disconnected from MongoDB")
	return nil
}
