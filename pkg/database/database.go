package database

import "go.mongodb.org/mongo-driver/mongo"

type Database interface {
	GetDb() *mongo.Database
}
