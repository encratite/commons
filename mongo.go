package commons

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func getCollection(collectionName, databaseName, uri string) *mongo.Collection {
	opts := options.Client()
	opts = opts.ApplyURI(uri)
	client, err := mongo.Connect(opts)
	if err != nil {
		Fatalf("Failed to connect to MongoDB server: %v", err)
	}
	database := client.Database(databaseName)
	collection := database.Collection(collectionName)
	return collection
}

func ReadMongo[T any](collectionName, databaseName, uri string) []T {
	collection := getCollection(collectionName, databaseName, uri)
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		Fatalf("Failed to find collection: %v", err)
	}
	var documents []T
	err = cursor.All(context.Background(), &documents)
	if err != nil {
		Fatalf("Failed to read collection: %v", err)
	}
	return documents
}

func WriteMongo[T any](collectionName, databaseName, uri string, documents []T) {
	start := time.Now()
	collection := getCollection(collectionName, databaseName, uri)
	collection.Drop(context.Background())
	manyOpts := options.InsertMany().SetOrdered(false)
	_, insertErr := collection.InsertMany(context.Background(), documents, manyOpts)
	if insertErr != nil {
		bulkErr, ok := insertErr.(mongo.BulkWriteException)
		if ok {
			for _, writeError := range bulkErr.WriteErrors {
				if writeError.Code != 11000 {
					Fatalf("Failed to insert document: %v", insertErr)
				}
			}
		}
	}
	end := time.Now()
	duration := end.Sub(start)
	fmt.Printf("Wrote %d documents to MongoDB in %d ms\n", len(documents), duration.Milliseconds())
}