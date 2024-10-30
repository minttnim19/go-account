package repository

import (
	"context"
	"go-account/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type oAuthRefreshTokenRepository struct {
	collection *mongo.Collection
}

func (r *oAuthRefreshTokenRepository) Create(token *model.OAuthRefreshToken) (*mongo.InsertOneResult, error) {
	token.Deleted = false
	token.CreatedAt = time.Now().Unix()
	return r.collection.InsertOne(context.TODO(), token)
}

func (r *oAuthRefreshTokenRepository) FindByID(id primitive.Binary) (model.OAuthRefreshToken, error) {
	token := model.OAuthRefreshToken{}
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id, "deleted": false}).Decode(&token)
	return token, err
}

func (r *oAuthRefreshTokenRepository) Update(id primitive.Binary, token *model.UpdateOAuthRefreshToken) error {
	token.UpdatedAt = time.Now().Unix()
	_, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": id, "deleted": false}, bson.M{"$set": token})
	return err
}

func NewOAuthRefreshTokenRepository(db *mongo.Database) model.OAuthRefreshTokenRepository {
	return &oAuthRefreshTokenRepository{collection: db.Collection("oauth_refresh_tokens")}
}
