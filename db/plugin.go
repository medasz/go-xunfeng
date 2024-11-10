package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-xunfeng/models"
)

func GetPluginAll(selector bson.M) ([]models.Plugin, error) {
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

func CountPluginAll(selecter bson.M) (int, error) {
	total, err := mgoDb.Collection("Plugin").CountDocuments(context.Background(), selecter)
	return int(total), err
}

func GetPlugin(selector bson.M, page, pageSize int) ([]models.Plugin, error) {
	result := make([]models.Plugin, 0)
	findOptions := new(options.FindOptions)
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64((page - 1) * pageSize))
	cur, err := mgoDb.Collection("Plugin").Find(context.Background(), selector, findOptions)
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

func GetVulCount() (int, error) {
	cur, err := mgoDb.Collection("Plugin").Aggregate(context.Background(), bson.A{
		bson.D{
			{
				"$group", bson.D{
					{"_id", ""},
					{"count", bson.D{
						{"$sum", "$count"},
					},
					},
				},
			},
		},
	})
	if err != nil {
		return 0, err
	}
	var list [][]bson.E
	err = cur.All(context.Background(), &list)
	if err != nil {
		return 0, err
	}
	if len(list) < 1 || len(list[0]) < 2 {
		return 0, nil
	}
	return int(list[0][1].Value.(int32)), nil
}

func GetVulGroupType() ([]models.VulType, error) {
	var data []models.VulType
	matchStage := bson.D{{"$match", bson.D{{"count", bson.D{{"$ne", 0}}}}}}
	groupStage := bson.D{{
		"$group", bson.D{
			{"_id", "$type"},
			{"count", bson.D{
				{"$sum", "$count"},
			},
			},
		},
	}}
	cur, err := mgoDb.Collection("Plugin").Aggregate(context.Background(), mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return data, err
	}
	var list [][]bson.E
	err = cur.All(context.Background(), &list)
	if err != nil {
		return data, err
	}
	for _, v := range list {
		if len(v) < 2 {
			continue
		}
		data = append(data, models.VulType{
			Type:  v[0].Value.(string),
			Count: int(v[1].Value.(int32)),
		})
	}
	return data, nil
}

func CreatePlugin(plugin models.Plugin) (*mongo.InsertOneResult, error) {
	return mgoDb.Collection("Plugin").InsertOne(context.Background(), plugin)
}
