package repositories

import (
	"context"
	"time"

	"go-account/internal/api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OAuthClientRepository interface {
	Create(user *models.CreateOAuthClient) (*mongo.InsertOneResult, error)
	Lists(filter map[string]interface{}, skip int, size int) ([]models.OAuthClient, int64, error)
	FindByID(id primitive.ObjectID) (models.OAuthClient, error)
}

type oAuthClientRepository struct {
	collection *mongo.Collection
}

func (r *oAuthClientRepository) Create(client *models.CreateOAuthClient) (*mongo.InsertOneResult, error) {
	client.Deleted = false
	client.CreatedAt = time.Now().Unix()
	return r.collection.InsertOne(context.TODO(), client)
}

func (r *oAuthClientRepository) FindByID(id primitive.ObjectID) (models.OAuthClient, error) {
	client := models.OAuthClient{}
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id, "deleted": false}).Decode(&client)
	return client, err
}

func (r *oAuthClientRepository) Lists(filter map[string]interface{}, skip int, size int) ([]models.OAuthClient, int64, error) {
	filter["deleted"] = false

	count, err := r.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, 0, err
	}
	findOptions := options.Find().SetSkip(int64(skip)).SetLimit(int64(size))
	cursor, err := r.collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, 0, err

	}
	clients := []models.OAuthClient{}
	err = cursor.All(context.TODO(), &clients)
	return clients, count, err
}

func NewOAuthClientRepository(db *mongo.Database) OAuthClientRepository {
	return &oAuthClientRepository{collection: db.Collection("oauth_clients")}
}
