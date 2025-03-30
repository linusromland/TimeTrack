package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var AuthDatabase *mongo.Database

func ConnectDB(mongoURI string) *mongo.Client {
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB")
	MongoClient = client
	AuthDatabase = client.Database("auth_service")
	return client
}

func GetCollection(collectionName string) *mongo.Collection {
	return AuthDatabase.Collection(collectionName)
}

func DisconnectDB() {
	if MongoClient != nil {
		if err := MongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
		log.Println("Disconnected from MongoDB")
	}
}