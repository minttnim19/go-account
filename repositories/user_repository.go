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

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUsers(filter map[string]interface{}, skip int, size int) ([]models.User, int64, error)
	FindUserByID(id primitive.ObjectID) (models.User, error)
	UpdateUser(id primitive.ObjectID, user *models.UpdateUser) error
	DeleteUser(id primitive.ObjectID) error
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{collection: db.Collection("users")}
}

func (r *userRepository) CreateUser(user *models.User) error {
	user.Deleted = false
	user.CreatedAt = time.Now().Unix()
	_, err := r.collection.InsertOne(context.TODO(), user)
	return err
}

func (r *userRepository) FindUserByID(id primitive.ObjectID) (models.User, error) {
	user := models.User{}
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id, "deleted": false}).Decode(&user)
	return user, err
}

func (r *userRepository) UpdateUser(id primitive.ObjectID, user *models.UpdateUser) error {
	user.UpdatedAt = time.Now().Unix()
	_, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": id, "deleted": false}, bson.M{"$set": user})
	return err
}

func (r *userRepository) DeleteUser(id primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"deleted":   true,
			"deletedAt": time.Now().Unix(),
		},
	}
	_, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	return err
}

func (r *userRepository) GetUsers(filter map[string]interface{}, skip int, size int) ([]models.User, int64, error) {
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
	users := []models.User{}
	err = cursor.All(context.TODO(), &users)
	return users, count, err
}
