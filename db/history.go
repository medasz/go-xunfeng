package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"go-xunfeng/models"
)

func CreateHistory(history models.History) (*mongo.InsertOneResult, error) {
	return mgoDb.Collection("History").InsertOne(context.Background(), history)
}
