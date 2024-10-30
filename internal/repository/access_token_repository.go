package repository

import (
	"context"
	"time"

	"go-account/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type oAuthAccessTokenRepository struct {
	collection *mongo.Collection
}

func (r *oAuthAccessTokenRepository) Create(token *model.OAuthAccessToken) (*mongo.InsertOneResult, error) {
	token.Deleted = false
	token.CreatedAt = time.Now().Unix()
	return r.collection.InsertOne(context.TODO(), token)
}

func (r *oAuthAccessTokenRepository) FindByID(id primitive.Binary) (model.OAuthAccessToken, error) {
	token := model.OAuthAccessToken{}
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id, "deleted": false}).Decode(&token)
	return token, err
}

func (r *oAuthAccessTokenRepository) Update(id primitive.Binary, token *model.UpdateOAuthAccessToken) error {
	token.UpdatedAt = time.Now().Unix()
	_, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": id, "deleted": false}, bson.M{"$set": token})
	return err
}

func NewOAuthAccessTokenRepository(db *mongo.Database) model.OAuthAccessTokenRepository {
	return &oAuthAccessTokenRepository{collection: db.Collection("oauth_access_tokens")}
}
