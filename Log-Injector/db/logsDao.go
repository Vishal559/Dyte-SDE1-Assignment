package db

import (
	"Log-Injector/controllers/logController"
	"Log-Injector/models"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database interface for handling database operations
type Database interface {
	InitLogEntryCollection() error
	InsertLogs(logs []models.LogEntry) error
	SearchLogsWithFilters(query string, filters map[string]string, page int, pageSize int) (string, error)
}

// MongoDB struct implements the Database interface
type MongoDB struct {
	logEntryCollection *mongo.Collection
}

// NewMongoDB creates a new MongoDB instance
func NewMongoDB(client *mongo.Client) *MongoDB {
	return &MongoDB{
		logEntryCollection: client.Database("LogInjector").Collection("logs"),
	}
}

// InitLogEntryCollection initializes the MongoDB collection for logs.
func (m *MongoDB) InitLogEntryCollection() error {
	// Get the MongoDB client
	client = GetClient()
	// Create index for the "Level" field
	indexModelLevel := mongo.IndexModel{
		Keys:    bson.D{{Key: "level", Value: 1}},
		Options: options.Index().SetBackground(true).SetSparse(true),
	}

	// Create index for the "Timestamp" field
	indexModelTimestamp := mongo.IndexModel{
		Keys:    bson.D{{Key: "timestamp", Value: 1}},
		Options: options.Index().SetBackground(true).SetSparse(true),
	}

	// Create index for the "Message" field
	indexModelMessage := mongo.IndexModel{
		Keys:    bson.D{{Key: "message", Value: "text"}},
		Options: options.Index().SetBackground(true).SetSparse(true),
	}

	// Create the indexes
	_, err := m.logEntryCollection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		indexModelLevel,
		indexModelTimestamp,
		indexModelMessage,
	})

	if err != nil {
		return fmt.Errorf("error creating indexes: %v", err)
	}

	fmt.Println("Indexes created for Level, Timestamp, and Message fields.")
	fmt.Println("Log entry collection initialized with indexes.")
	return nil
}

// InsertLogs inserts multiple log entries into the MongoDB collection.
func (m *MongoDB) InsertLogs(logEntries []models.LogEntry) error {

	// Convert LogEntry objects to BSON documents
	var documents []interface{}
	for _, entry := range logEntries {
		documents = append(documents, entry)
	}

	// Perform the bulk insert
	result, err := m.logEntryCollection.InsertMany(context.TODO(), documents)
	if err != nil {
		return fmt.Errorf("error inserting log entries: %v", err)
	}

	fmt.Printf("Inserted %v documents\n", len(result.InsertedIDs))
	return nil
}

// SearchLogsWithFilters performs a text-based search with multiple filters on specified fields.
func (m *MongoDB) SearchLogsWithFilters(query string, filters map[string]string, page int, pageSize int) (string, error) {
	// Create a case-insensitive regular expression for the query
	regex := bson.M{"$regex": primitive.Regex{Pattern: query, Options: "i"}}

	// Create a map to hold the filters
	filterMap := bson.M{"message": regex}

	// Add additional filters based on the provided filter map
	for field, value := range filters {
		if field == "parentResourceId" {
			filterMap["metadata"] = bson.M{field: value}
		} else {
			filterMap[field] = value
		}
	}

	// Perform the text search with filters and pagination
	cur, err := m.logEntryCollection.Find(context.TODO(), filterMap, options.Find().SetLimit(int64(pageSize)).SetSkip(int64((page-1)*pageSize)))
	if err != nil {
		return "", fmt.Errorf("error performing text search: %v", err)
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err = cur.Close(ctx)
		if err != nil {
			return
		}
	}(cur, context.TODO())

	// Decode the results into LogEntry objects
	var results []logController.LogEntry
	for cur.Next(context.TODO()) {
		var logEntry logController.LogEntry
		err = cur.Decode(&logEntry)
		if err != nil {
			return "", fmt.Errorf("error decoding log entry: %v", err)
		}
		results = append(results, logEntry)
	}

	// Build JSON response
	jsonResponse, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error encoding JSON: %v", err)
	}

	return string(jsonResponse), nil
}
