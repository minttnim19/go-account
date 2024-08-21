package repositories

import (
	"context"
	"time"

	"go-account/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OAuthClientRepository interface {
	Create(user *models.CreateOAuthClient) (*mongo.InsertOneResult, error)
	Lists(filter map[string]interface{}, skip int, size int) ([]models.OAuthClient, int64, error)
	// FindUserByUsername(username string) (models.User, error)
	FindByID(id primitive.ObjectID) (models.OAuthClient, error)
	// UpdateUser(id primitive.ObjectID, user *models.UpdateUser) error
	// DeleteUser(id primitive.ObjectID) error
}

type oAuthClientRepository struct {
	collection *mongo.Collection
}

func (r *oAuthClientRepository) Create(client *models.CreateOAuthClient) (*mongo.InsertOneResult, error) {
	client.Deleted = false
	client.CreatedAt = time.Now().Unix()
	return r.collection.InsertOne(context.TODO(), client)
}

// func (r *userRepository) FindUserByUsername(username string) (models.User, error) {
// 	user := models.User{}
// 	err := r.collection.FindOne(context.TODO(), bson.M{"username": username, "deleted": false}).Decode(&user)
// 	return user, err
// }

func (r *oAuthClientRepository) FindByID(id primitive.ObjectID) (models.OAuthClient, error) {
	client := models.OAuthClient{}
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id, "deleted": false}).Decode(&client)
	return client, err
}

// func (r *userRepository) UpdateUser(id primitive.ObjectID, user *models.UpdateUser) error {
// 	user.UpdatedAt = time.Now().Unix()
// 	_, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": id, "deleted": false}, bson.M{"$set": user})
// 	return err
// }

// func (r *userRepository) DeleteUser(id primitive.ObjectID) error {
// 	update := bson.M{
// 		"$set": bson.M{
// 			"deleted":   true,
// 			"deletedAt": time.Now().Unix(),
// 		},
// 	}
// 	_, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
// 	return err
// }

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
