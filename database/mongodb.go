package database

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var databaseName string

func InitMongoDB() error {
	mongoUri := os.Getenv("MONGO")
	databaseName = os.Getenv("DATABASE")

	var mongoClientErr error
	mongoClient, mongoClientErr = mongo.Connect(context.Background(), options.Client().ApplyURI(mongoUri))
	if mongoClientErr != nil {
		panic(mongoClientErr)
	}

	fmt.Println("mongo client", mongoClient)

	pingErr := mongoClient.Ping(context.Background(), nil)
	fmt.Println("ping", pingErr)
	if pingErr != nil {
		return errors.New("connection cannto be made")
	}

	return nil
}

func DisconnectMongoDB() {
	fmt.Println("disconnect")
	err := mongoClient.Disconnect(context.Background())
	if err != nil {
		panic(err)
	}
}

func FetchCollection(name string) *mongo.Collection {
	return mongoClient.Database("cpu-mon").Collection("users")
}
