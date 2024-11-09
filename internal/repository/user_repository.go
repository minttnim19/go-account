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

type UserRepository interface {
	Create(user *domain.User) error
	Lists(filter map[string]interface{}, skip int, size int) ([]domain.User, int64, error)
	FindUserByUsername(username string) (domain.User, error)
	FindByID(id primitive.ObjectID) (domain.User, error)
	Update(id primitive.ObjectID, user *domain.UpdateUser) error
	Delete(id primitive.ObjectID) error
}

type userRepository struct {
	collection *mongo.Collection
}

func (r *userRepository) Create(user *domain.User) error {
	user.Deleted = false
	user.CreatedAt = time.Now().Unix()
	_, err := r.collection.InsertOne(context.TODO(), user)
	return err
}

func (r *userRepository) FindUserByUsername(username string) (domain.User, error) {
	user := domain.User{}
	err := r.collection.FindOne(context.TODO(), bson.M{"username": username, "deleted": false}).Decode(&user)
	return user, err
}

func (r *userRepository) FindByID(id primitive.ObjectID) (domain.User, error) {
	user := domain.User{}
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id, "deleted": false}).Decode(&user)
	return user, err
}

func (r *userRepository) Update(id primitive.ObjectID, user *domain.UpdateUser) error {
	updateFields := bson.M{"updatedAt": time.Now().Unix()}
	fieldMapping := map[string]interface{}{
		"status": user.Status,
	}
	utils.FieldMapping(fieldMapping, &updateFields)
	_, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": id, "deleted": false}, bson.M{"$set": updateFields})
	return err
}

func (r *userRepository) Delete(id primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"deleted":   true,
			"deletedAt": time.Now().Unix(),
		},
	}
	_, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	return err
}

func (r *userRepository) Lists(filter map[string]interface{}, skip int, size int) ([]domain.User, int64, error) {
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
	users := []domain.User{}
	err = cursor.All(context.TODO(), &users)
	return users, count, err
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{collection: db.Collection("users")}
}
