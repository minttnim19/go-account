package repository

import (
	"context"
	"time"

	"go-account/internal/domain"
	"go-account/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OAuthClientRepository interface {
	Create(user *domain.OAuthClient) (*mongo.InsertOneResult, error)
	Lists(filter map[string]interface{}, skip int, size int) ([]domain.OAuthClient, int64, error)
	FindByID(id primitive.ObjectID) (domain.OAuthClient, error)
	Update(id primitive.ObjectID, client *domain.UpdateOAuthClient) error
	Delete(id primitive.ObjectID) error
}
type oAuthClientRepository struct {
	collection *mongo.Collection
}

func (r *oAuthClientRepository) Create(client *domain.OAuthClient) (*mongo.InsertOneResult, error) {
	client.Deleted = false
	client.CreatedAt = time.Now().Unix()
	return r.collection.InsertOne(context.TODO(), client)
}

func (r *oAuthClientRepository) FindByID(id primitive.ObjectID) (domain.OAuthClient, error) {
	client := domain.OAuthClient{}
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id, "deleted": false}).Decode(&client)
	return client, err
}

func (r *oAuthClientRepository) Update(id primitive.ObjectID, client *domain.UpdateOAuthClient) error {
	updateFields := bson.M{"updatedAt": time.Now().Unix()}
	fieldMapping := map[string]interface{}{
		"name":       client.Name,
		"redirects":  client.Redirects,
		"scopes":     client.Scopes,
		"grantTypes": client.GrantTypes,
		"revoked":    client.Revoked,
	}
	utils.FieldMapping(fieldMapping, &updateFields)

	_, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": id, "deleted": false}, bson.M{"$set": updateFields})
	return err
}

func (r *oAuthClientRepository) Delete(id primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"deleted":   true,
			"deletedAt": time.Now().Unix(),
		},
	}
	_, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	return err
}

func (r *oAuthClientRepository) Lists(filter map[string]interface{}, skip int, size int) ([]domain.OAuthClient, int64, error) {
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
	clients := []domain.OAuthClient{}
	err = cursor.All(context.TODO(), &clients)
	return clients, count, err
}

func NewOAuthClientRepository(db *mongo.Database) OAuthClientRepository {
	return &oAuthClientRepository{collection: db.Collection("oauth_clients")}
}
