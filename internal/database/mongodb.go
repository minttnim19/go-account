package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongo(mongoDBURI string, mongoDBUser string, mongoDBPassword string, mongoDBName string) *mongo.Database {
	// Replace placeholders in MONGO_URI with actual credentials
	uri := fmt.Sprintf(mongoDBURI, mongoDBUser, mongoDBPassword)

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
	database := client.Database(mongoDBName)

	createIndex(database) // Create mongodb index
	return database
}

func createIndex(database *mongo.Database) {
	collUser := database.Collection("users")
	userIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: -1}},
		Options: options.Index().SetUnique(true), // unique index
	}
	if _, err := collUser.Indexes().CreateOne(context.TODO(), userIndex); err != nil {
		panic(err)
	}

	collOAuthClient := database.Collection("oauth_clients")
	oAuthClientIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "_id", Value: -1}, {Key: "secret", Value: -1}},
		Options: options.Index().SetUnique(true), // unique index
	}
	if _, err := collOAuthClient.Indexes().CreateOne(context.TODO(), oAuthClientIndex); err != nil {
		panic(err)
	}

	collOAuthRefreshToken := database.Collection("oauth_refresh_tokens")
	oAuthRefreshTokenIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "accessTokenID", Value: -1}},
		Options: options.Index(), // unique index
	}
	if _, err := collOAuthRefreshToken.Indexes().CreateOne(context.TODO(), oAuthRefreshTokenIndex); err != nil {
		panic(err)
	}
}
