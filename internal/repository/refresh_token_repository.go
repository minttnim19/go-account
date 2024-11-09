package repository

import (
	"context"
	"go-account/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OAuthRefreshTokenRepository interface {
	Create(token *domain.OAuthRefreshToken) (*mongo.InsertOneResult, error)
	FindByID(id primitive.Binary) (domain.OAuthRefreshToken, error)
	Update(id primitive.Binary, token *domain.UpdateOAuthRefreshToken) error
}

type oAuthRefreshTokenRepository struct {
	collection *mongo.Collection
}

func (r *oAuthRefreshTokenRepository) Create(token *domain.OAuthRefreshToken) (*mongo.InsertOneResult, error) {
	token.Deleted = false
	token.CreatedAt = time.Now().Unix()
	return r.collection.InsertOne(context.TODO(), token)
}

func (r *oAuthRefreshTokenRepository) FindByID(id primitive.Binary) (domain.OAuthRefreshToken, error) {
	token := domain.OAuthRefreshToken{}
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id, "deleted": false}).Decode(&token)
	return token, err
}

func (r *oAuthRefreshTokenRepository) Update(id primitive.Binary, token *domain.UpdateOAuthRefreshToken) error {
	token.UpdatedAt = time.Now().Unix()
	_, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": id, "deleted": false}, bson.M{"$set": token})
	return err
}

func NewOAuthRefreshTokenRepository(db *mongo.Database) OAuthRefreshTokenRepository {
	return &oAuthRefreshTokenRepository{collection: db.Collection("oauth_refresh_tokens")}
}
