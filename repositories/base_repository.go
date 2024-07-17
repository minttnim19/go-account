package repositories

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BaseRepository struct {
	collection *mongo.Collection
}

func structToMap(data interface{}) (map[string]interface{}, error) {
	dataType := reflect.TypeOf(data)

	if dataType.Kind() == reflect.Ptr {
		dataType = dataType.Elem()
		data = reflect.ValueOf(data).Elem().Interface()
	}

	if dataType.Kind() != reflect.Struct {
		return nil, errors.New("input data is not a struct")
	}

	dataValue := reflect.ValueOf(data)
	dataMap := make(map[string]interface{})
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		fieldValue := dataValue.Field(i)
		fieldName := field.Name
		lowerFirst := strings.ToLower(fieldName[:1]) + fieldName[1:]
		dataMap[lowerFirst] = fieldValue.Interface()
	}

	return dataMap, nil
}

func (b *BaseRepository) Create(document interface{}) error {
	docMap, err := structToMap(document)
	if err != nil {
		return err
	}
	docMap["deleted"] = false
	docMap["createdAt"] = time.Now().Unix()
	_, err = b.collection.InsertOne(context.TODO(), docMap)
	return err
}

func (b *BaseRepository) FindByID(id primitive.ObjectID) *mongo.SingleResult {
	return b.collection.FindOne(context.Background(), bson.M{"_id": id, "deleted": false})
}

func (b *BaseRepository) Get(filter map[string]interface{}) (*mongo.Cursor, error) {
	filter["deleted"] = false
	return b.collection.Find(context.Background(), filter)
}
