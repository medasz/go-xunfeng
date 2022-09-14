package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-xunfeng/models"
)

func GetTask(selector bson.M, page int, pageSize int) ([]models.Task, error) {
	result := make([]models.Task, 0)
	findOptions := new(options.FindOptions)
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64((page - 1) * pageSize))
	findOptions.SetSort(bson.M{"time": -1})
	cur, err := mgoDb.Collection("Task").Find(context.Background(), selector, findOptions)
	if err != nil {
		return result, err
	}
	for cur.Next(context.Background()) {
		tmp := models.Task{}
		err := cur.Decode(&tmp)
		if err != nil {
			return result, err
		}
		result = append(result, tmp)
	}

	return result, nil
}

func CountAllTask(selector bson.M) (int, error) {
	total, err := mgoDb.Collection("Task").CountDocuments(context.Background(), selector)
	return int(total), err

}

func TaskDeleteAll() error {
	_, err := mgoDb.Collection("Task").DeleteMany(context.Background(), bson.M{})
	return err
}

func CreateTask(task models.InTask) (*mongo.InsertOneResult, error) {
	return mgoDb.Collection("Task").InsertOne(context.Background(), task)
}

func DeleteTaskById(id string) (*mongo.DeleteResult, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return mgoDb.Collection("Task").DeleteOne(context.Background(), bson.M{"_id": objectId})
}

func GetTaskAll() ([]models.Task, error) {
	result := make([]models.Task, 0)
	cur, err := mgoDb.Collection("Task").Find(context.Background(), bson.M{})
	if err != nil {
		return result, err
	}
	for cur.Next(context.Background()) {
		tmp := models.Task{}
		err := cur.Decode(&tmp)
		if err != nil {
			return result, err
		}
		result = append(result, tmp)
	}
	return result, nil
}
