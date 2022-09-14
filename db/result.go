package db

import (
	"context"

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
