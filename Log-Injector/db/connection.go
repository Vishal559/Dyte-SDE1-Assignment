package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://admin-vishal:RSjrcspwOLrYs7Ex@cluster0.cujjf.mongodb.net/?retryWrites=true&w=majority"

var client *mongo.Client

func Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(connectionString)
	c, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	// Check the connection
	err = c.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("error pinging MongoDB: %v", err)
	}

	client = c
	fmt.Println("Connected to MongoDB!")

	return nil
}

// GetClient returns the MongoDB client.
func GetClient() *mongo.Client {
	return client
}

// GetContext returns a context for MongoDB operations.
func GetContext() context.Context {
	return context.TODO()
}
