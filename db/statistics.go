package db

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-xunfeng/models"
)

func GetStatistics() (models.Statistics, error) {
	nowDay := time.Now().Format("2006-01-02")
	data := models.Statistics{}
	selecter := bson.M{
		"date": nowDay,
	}
	err := mgoDb.Collection("Statistics").FindOne(context.Background(), selecter).Decode(&data)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return data, nil
	}

	return data, err
}

func GetAllStatistics() ([]models.Statistics, error) {
	data := make([]models.Statistics, 0)
	cur, err := mgoDb.Collection("Statistics").Find(context.Background(), bson.M{})
	if err != nil {
		return data, err
	}
	defer cur.Close(context.Background())
	err = cur.All(context.Background(), &data)
	return data, err
}

func UpdateOrUpsertStatistic(date string, info models.StatisticsInfo) error {
	_, err := mgoDb.Collection("Statistics").UpdateOne(context.Background(),
		bson.M{"date": date}, bson.M{"$set": bson.M{"info": info}}, options.Update().SetUpsert(true))
	return err
}

func GetStatisticsLimit() ([]models.Statistics, error) {
	result := make([]models.Statistics, 0)
	findOptions := new(options.FindOptions)
	findOptions.SetLimit(30)
	findOptions.SetSort(bson.M{"date": -1})
	cur, err := mgoDb.Collection("Statistics").Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		return result, err
	}
	for cur.Next(context.Background()) {
		tmp := models.Statistics{}
		err := cur.Decode(&tmp)
		if err != nil {
			return result, err
		}
		result = append(result, tmp)
	}

	return result, nil
}
