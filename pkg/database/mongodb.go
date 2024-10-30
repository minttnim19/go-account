package database

import (
	"context"
	"fmt"
	"go-account/config"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDatabase struct {
	Db *mongo.Database
}

var (
	once       sync.Once
	dbInstance *mongoDatabase
)

func NewMongoDatabase(conf *config.Config) Database {
	once.Do(func() {
		uri := fmt.Sprintf(conf.MongoDBURI, conf.MongoDBUser, conf.MongoDBPassword)

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
		database := client.Database(conf.MongoDBName)

		// Ensure Indexes
		ensureIndexes(database)

		dbInstance = &mongoDatabase{Db: database}
	})

	return dbInstance
}

func ensureIndexes(database *mongo.Database) {
	createIndex := func(coll *mongo.Collection, index mongo.IndexModel) {
		if _, err := coll.Indexes().CreateOne(context.TODO(), index); err != nil {
			panic(err)
		}
	}

	// Define indexes
	userIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: -1}},
		Options: options.Index().SetUnique(true), // unique index
	}
	oAuthClientIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "_id", Value: -1}, {Key: "secret", Value: -1}},
		Options: options.Index().SetUnique(true), // unique index
	}
	oAuthRefreshTokenIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "accessTokenID", Value: -1}},
		Options: options.Index(),
	}

	// Create indexes
	createIndex(database.Collection("users"), userIndex)
	createIndex(database.Collection("oauth_clients"), oAuthClientIndex)
	createIndex(database.Collection("oauth_refresh_tokens"), oAuthRefreshTokenIndex)
}

func (p *mongoDatabase) GetDb() *mongo.Database {
	return dbInstance.Db
}
