package db

import (
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go-xunfeng/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteResultByTaskId(taskId string) (*mongo.DeleteResult, error) {
	objectId, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		return nil, err
	}
	return mgoDb.Collection("Result").DeleteMany(context.Background(), bson.M{"task_id": objectId})
}

func GetDistinctTaskDateByTaskId(taskId string) ([]time.Time, error) {
	data := make([]time.Time, 0)
	objectId, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		return data, err
	}

	res, err := mgoDb.Collection("Result").Distinct(context.Background(), "task_date", bson.M{"task_id": objectId})
	if err != nil {
		return data, err
	}
	tmp, err := json.Marshal(res)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(tmp, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func GetResult(selector bson.M, optionList *options.FindOptions) ([]models.Result, error) {
	result := make([]models.Result, 0)
	cur, err := mgoDb.Collection("Result").Find(context.Background(), selector, optionList)
	if err != nil {
		return result, err
	}
	defer cur.Close(context.Background())
	err = cur.All(context.Background(), &result)
	return result, err
}
