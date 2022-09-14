package db

import (
	"context"

	"go-xunfeng/models"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdateHeartbeat(data bson.M) error {
	_, err := mgoDb.Collection("Heartbeat").UpdateOne(context.Background(), bson.M{"name": "heartbeat"}, data)
	return err
}

func GetHeartbeat(name string) (models.Heartbeat, error) {
	data := models.Heartbeat{}
	err := mgoDb.Collection("Heartbeat").FindOne(context.Background(), bson.M{"name": name}).Decode(&data)
	return data, err
}
