package db

import (
	"context"
	"go-log-keeper/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func NewMongoClient() (*mongo.Client, *mongo.Collection) {
	clientOpts := options.Client().ApplyURI(config.MongoURI)
	mongoClient, err := mongo.Connect(context.Background(), clientOpts)

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v -- %s", err, config.MongoURI)
	}
	collection := mongoClient.Database(config.MongoDB).Collection(config.MongoCollection)
	return mongoClient, collection
}
