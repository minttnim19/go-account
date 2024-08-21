package repositories

import (
	"context"
	"time"

	"go-account/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OAuthAccessTokenRepository interface {
	Create(token *models.OAuthAccessToken) (*mongo.InsertOneResult, error)
	CreateWithSession(sessionCtx mongo.SessionContext, token *models.OAuthAccessToken) (*mongo.InsertOneResult, error)
	FindByID(id primitive.ObjectID) (models.OAuthAccessToken, error)
	Update(id primitive.ObjectID, user *models.UpdateOAuthAccessToken) error
}

type oAuthAccessTokenRepository struct {
	collection *mongo.Collection
}

func (r *oAuthAccessTokenRepository) Create(token *models.OAuthAccessToken) (*mongo.InsertOneResult, error) {
	token.Deleted = false
	token.CreatedAt = time.Now().Unix()
	return r.collection.InsertOne(context.TODO(), token)
}

func (r *oAuthAccessTokenRepository) CreateWithSession(sessionCtx mongo.SessionContext, token *models.OAuthAccessToken) (*mongo.InsertOneResult, error) {
	token.Deleted = false
	token.CreatedAt = time.Now().Unix()
	return r.collection.InsertOne(sessionCtx, token)
}

func (r *oAuthAccessTokenRepository) FindByID(id primitive.ObjectID) (models.OAuthAccessToken, error) {
	token := models.OAuthAccessToken{}
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id, "deleted": false}).Decode(&token)
	return token, err
}

func (r *oAuthAccessTokenRepository) Update(id primitive.ObjectID, token *models.UpdateOAuthAccessToken) error {
	token.UpdatedAt = time.Now().Unix()
	_, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": id, "deleted": false}, bson.M{"$set": token})
	return err
}

func NewOAuthAccessTokenRepository(db *mongo.Database) OAuthAccessTokenRepository {
	return &oAuthAccessTokenRepository{collection: db.Collection("oauth_access_tokens")}
}
