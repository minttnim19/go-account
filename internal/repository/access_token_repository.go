package repository

import (
	"context"
	"time"

	"go-account/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OAuthAccessTokenRepository interface {
	Create(token *domain.OAuthAccessToken) (*mongo.InsertOneResult, error)
	FindByID(id primitive.Binary) (domain.OAuthAccessToken, error)
	Update(id primitive.Binary, user *domain.UpdateOAuthAccessToken) error
}

type oAuthAccessTokenRepository struct {
	collection *mongo.Collection
}

func (r *oAuthAccessTokenRepository) Create(token *domain.OAuthAccessToken) (*mongo.InsertOneResult, error) {
	token.Deleted = false
	token.CreatedAt = time.Now().Unix()
	return r.collection.InsertOne(context.TODO(), token)
}

func (r *oAuthAccessTokenRepository) FindByID(id primitive.Binary) (domain.OAuthAccessToken, error) {
	token := domain.OAuthAccessToken{}
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id, "deleted": false}).Decode(&token)
	return token, err
}

func (r *oAuthAccessTokenRepository) Update(id primitive.Binary, token *domain.UpdateOAuthAccessToken) error {
	token.UpdatedAt = time.Now().Unix()
	_, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": id, "deleted": false}, bson.M{"$set": token})
	return err
}

func NewOAuthAccessTokenRepository(db *mongo.Database) OAuthAccessTokenRepository {
	return &oAuthAccessTokenRepository{collection: db.Collection("oauth_access_tokens")}
}
