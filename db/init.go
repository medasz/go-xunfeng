package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-xunfeng/config"
)

var mgoCli *mongo.Client
var mgoDb *mongo.Database

func init() {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
		config.Cfg.MongoDb.DbUsername,
		config.Cfg.MongoDb.DbPassword,
		config.Cfg.MongoDb.Host,
		config.Cfg.MongoDb.Port,
		config.Cfg.MongoDb.DbName)
	clientOption := options.Client().ApplyURI(uri)
	mgoCli, err := mongo.Connect(context.Background(), clientOption)
	if err != nil {
		panic(err)
	}
	mgoDb = mgoCli.Database(config.Cfg.MongoDb.DbName)
}
