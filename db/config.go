package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-xunfeng/models"
)

func GetNascanConfig() (*models.Nascan, error) {
	result := new(models.Nascan)
	err := mgoDb.Collection("Config").FindOne(context.Background(), map[string]string{"type": "nascan"}, options.FindOne()).Decode(&result)
	return result, err
}
func GetVulscanConfig() (*models.Vulscan, error) {
	result := new(models.Vulscan)
	err := mgoDb.Collection("Config").FindOne(context.Background(), map[string]string{"type": "vulscan"}, options.FindOne()).Decode(&result)
	return result, err
}

func UpdateConfig(selector bson.M, data bson.M) error {
	_, err := mgoDb.Collection("Config").UpdateOne(context.Background(), selector, data)
	return err
}
