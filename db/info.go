package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-xunfeng/models"
)

func GetInfo(selector bson.M, page int, pageSize int) ([]models.Info, error) {
	data := make([]models.Info, 0)
	findOptions := new(options.FindOptions)
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64((page - 1) * pageSize))
	findOptions.SetSort(bson.M{"time": -1})
	cur, err := mgoDb.Collection("Info").Find(context.Background(), selector, findOptions)
	if err != nil {
		return data, err
	}
	defer cur.Close(context.Background())
	err = cur.All(context.Background(), &data)
	return data, err
}

func CountAll(selector bson.M) (int, error) {
	total, err := mgoDb.Collection("Info").CountDocuments(context.Background(), selector)
	return int(total), err

}

func GetInfoAllIpPort(selector bson.M) ([][]interface{}, error) {
	data := make([][]interface{}, 0)

	cur, err := mgoDb.Collection("Info").Find(context.Background(), selector)
	if err != nil {
		return data, err
	}
	defer cur.Close(context.Background())
	tmpData := make([]models.Info, 0)
	err = cur.All(context.Background(), &tmpData)
	if err != nil {
		return data, err
	}
	for _, tmp := range tmpData {
		data = append(data, []interface{}{tmp.Ip, tmp.Port})
	}
	return data, nil
}
