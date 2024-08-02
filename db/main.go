package db

import (
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"context"
)

var Client *mongo.Client

func CreateDBClient() {
	// Set up a MongoDB client
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	Client = client
}

func GetUsersCollection() *mongo.Collection {
	return Client.Database("trashauth").Collection("Users")
}