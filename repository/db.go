package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"fortknox-api/model"
)

const uri = "mongodb://127.0.0.1:27017"

var collection = GetMongoCollection()

func ConnectToDB() *mongo.Client {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Pinged your datastore. You have successfully connected to MongoDB!")
	return client
}

func GetMongoCollection() *mongo.Collection {
	return ConnectToDB().Database("fortknox-db").Collection("tokenStore")
}

func FindBy(filter bson.D, c chan func() (model.MongoDBTokenStoreDocument, error)) {
	// annonymous function made to return 2 values to channel
	c <- (func() (model.MongoDBTokenStoreDocument, error) {
		var result model.MongoDBTokenStoreDocument
		err := collection.FindOne(context.TODO(), filter).Decode(&result)
		log.Println("Processed from dataStore Channel")
		return result, err
	})
}
