package repositories

import (
	"context"
	"time"

	"go-account/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OAuthRefreshTokenRepository interface {
	Create(token *models.OAuthRefreshToken) (*mongo.InsertOneResult, error)
	CreateWithSession(sessionCtx mongo.SessionContext, token *models.OAuthRefreshToken) (*mongo.InsertOneResult, error)
	FindByID(id primitive.ObjectID) (models.OAuthRefreshToken, error)
	Update(id primitive.ObjectID, user *models.UpdateOAuthRefreshToken) error
}

type oAuthRefreshTokenRepository struct {
	collection *mongo.Collection
}

func (r *oAuthRefreshTokenRepository) Create(token *models.OAuthRefreshToken) (*mongo.InsertOneResult, error) {
	token.Deleted = false
	token.CreatedAt = time.Now().Unix()
	return r.collection.InsertOne(context.TODO(), token)
}

func (r *oAuthRefreshTokenRepository) CreateWithSession(sessionCtx mongo.SessionContext, token *models.OAuthRefreshToken) (*mongo.InsertOneResult, error) {
	token.Deleted = false
	token.CreatedAt = time.Now().Unix()
	return r.collection.InsertOne(sessionCtx, token)
}

func (r *oAuthRefreshTokenRepository) FindByID(id primitive.ObjectID) (models.OAuthRefreshToken, error) {
	token := models.OAuthRefreshToken{}
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id, "deleted": false}).Decode(&token)
	return token, err
}

func (r *oAuthRefreshTokenRepository) Update(id primitive.ObjectID, token *models.UpdateOAuthRefreshToken) error {
	token.UpdatedAt = time.Now().Unix()
	_, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": id, "deleted": false}, bson.M{"$set": token})
	return err
}

func NewOAuthRefreshTokenRepository(db *mongo.Database) OAuthRefreshTokenRepository {
	return &oAuthRefreshTokenRepository{collection: db.Collection("oauth_refresh_tokens")}
}
