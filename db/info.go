package db

import (
	"context"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func CountAllInfo(selector bson.M) (int, error) {
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

func GetOneInfoByIp(ip string) (*models.Info, error) {
	data := new(models.Info)
	err := mgoDb.Collection("Info").FindOne(context.Background(), bson.M{"ip": ip}).Decode(&data)
	return data, err
}

func GetInfoAll() ([]models.Info, error) {
	data := make([]models.Info, 0)
	findOptions := new(options.FindOptions)
	findOptions.SetSort(bson.M{"time": 1})
	cur, err := mgoDb.Collection("Info").Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		return data, err
	}
	defer cur.Close(context.Background())
	err = cur.All(context.Background(), &data)
	return data, err
}

func DeleteInfo(ip string, port int) (*mongo.DeleteResult, error) {
	return mgoDb.Collection("Info").DeleteOne(context.Background(), bson.M{"ip": ip, "port": port})
}

func GetInfoIpCount() (int, error) {
	ips, err := mgoDb.Collection("Info").Distinct(context.Background(), "ip", bson.M{})
	return len(ips), err
}

func GetServerType() ([]models.TypeCount, error) {
	var data []models.TypeCount
	matchStage := bson.D{{
		"$match", bson.D{
			{"server", bson.D{{"$ne", "web"}}},
		},
	}}
	groupStage := bson.D{{
		"$group", bson.D{
			{"_id", "$server"},
			{"count", bson.D{
				{"$sum", 1},
			},
			},
		},
	}}
	sortStage := bson.D{{
		"$sort", bson.D{{"count", -1}},
	}}
	cur, err := mgoDb.Collection("Info").Aggregate(context.Background(), mongo.Pipeline{matchStage, groupStage, sortStage})
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
		if v[0].Value != nil {
			data = append(data, models.TypeCount{
				ServerName: v[0].Value.(string),
				Count:      int(v[1].Value.(int32)),
			})
		}
	}
	return data, nil
}

func GetWebType() ([]models.TypeCount, error) {
	var data []models.TypeCount
	matchStage := bson.D{{
		"$match", bson.D{
			{"server", "web"},
		},
	}}
	unwindStage := bson.D{{
		"$unwind", "$webinfo.tag",
	},
	}
	groupStage := bson.D{{
		"$group", bson.D{
			{"_id", "$webinfo.tag"},
			{"count", bson.D{
				{"$sum", 1},
			},
			},
		},
	}}
	sortStage := bson.D{{
		"$sort", bson.D{{"count", -1}},
	}}
	cur, err := mgoDb.Collection("Info").Aggregate(context.Background(), mongo.Pipeline{matchStage, unwindStage, groupStage, sortStage})
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
		if v[0].Value != nil {
			data = append(data, models.TypeCount{
				ServerName: v[0].Value.(string),
				Count:      int(v[1].Value.(int32)),
			})
		}
	}
	return data, nil
}

func CreateInfo(info models.Info) (*mongo.InsertOneResult, error) {
	return mgoDb.Collection("Info").InsertOne(context.Background(), info)
}

func FindOneInfoAndDelete(selector bson.M) *mongo.SingleResult {
	return mgoDb.Collection("Info").FindOneAndDelete(context.Background(), selector, options.FindOneAndDelete())
}

func UpdateInfo(ip, port, server string) error {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return err
	}
	_, err = mgoDb.Collection("Info").UpdateOne(context.Background(), bson.M{"ip": ip, "port": portInt}, bson.M{"$set": bson.M{"server": server}})
	return err
}

func UpdateInfoAll(ip, port string, data bson.M) error {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return err
	}
	_, err = mgoDb.Collection("Info").UpdateOne(context.Background(), bson.M{"ip": ip, "port": portInt}, data)
	return err
}

func GetInfoOne(selector bson.M) (*models.Info, error) {
	data := new(models.Info)
	err := mgoDb.Collection("Info").FindOne(context.Background(), selector).Decode(&data)
	return data, err
}
