package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"go-xunfeng/models"
)

func GetPlugin(selector bson.M) ([]models.Plugin, error) {
	result := make([]models.Plugin, 0)
	cur, err := mgoDb.Collection("Plugin").Find(context.Background(), selector)
	if err != nil {
		return result, err
	}
	for cur.Next(context.Background()) {
		tmp := models.Plugin{}
		err := cur.Decode(&tmp)
		if err != nil {
			return result, err
		}
		result = append(result, tmp)
	}

	return result, nil
}

func GetPluginTypes() ([]string, error) {
	data := make([]string, 0)
	tmp, err := mgoDb.Collection("Plugin").Distinct(context.Background(), "type", bson.M{})
	for _, v := range tmp {
		data = append(data, v.(string))
	}
	return data, err
}
