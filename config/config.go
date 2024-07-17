package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect2MongoDB() (*mongo.Database, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Environment variable not set")
	}

	mongoURI := os.Getenv("MONGO_URI")
	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	databaseName := os.Getenv("DATABASE_NAME")

	// Replace placeholders in MONGO_URI with actual credentials
	uri := fmt.Sprintf(mongoURI, mongoUser, mongoPassword)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	database := client.Database(databaseName)

	createIndex(database) // Create mongodb index
	return database, err
}

func createIndex(database *mongo.Database) {
	collUser := database.Collection("users")
	_, err := collUser.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true), // unique index
	})
	if err != nil {
		panic(err)
	}
}
