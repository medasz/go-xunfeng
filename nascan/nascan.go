package nascan

import (
	"encoding/base64"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"go-xunfeng/pkg/nascan"

	"go.mongodb.org/mongo-driver/bson"

	"go-xunfeng/models"

	"go-xunfeng/db"
)

var (
	masscan  int64
	nachange int64
	//statistics = make(map[string]models.StatisticsInfo)
	//configInfo *models.Nascan
	statisticsAtomic atomic.Value
	configInfoAtomic atomic.Value
	err              error
)

func init() {
	configInfo, err := nascan.GetNascanConfig()
	if err != nil {
		panic(err)
	}
	log.Println("获取配置成功")
	configInfoAtomic.Store(configInfo)
	statistic, err := db.GetStatistics()
	if err != nil {
		panic(err)
	}
	statistics := make(map[string]models.StatisticsInfo)
	statistics[time.Now().Format("2006-01-02")] = statistic.Info
	statisticsAtomic.Store(statistics)
}

func monitor() {
	for true {
		err = db.UpdateHeartbeat(bson.M{"$set": bson.M{"up_time": time.Now()}})
		if err != nil {
			log.Println(err)
		}
		statistics := statisticsAtomic.Load().(map[string]models.StatisticsInfo)
		if _, ok := statistics[time.Now().Format("2006-01-02")]; !ok {
			statistics[time.Now().Format("2006-01-02")] = models.StatisticsInfo{
				Add:    0,
				Update: 0,
				Delete: 0,
			}
		}
		err := db.UpdateOrUpsertStatistic(time.Now().Format("2006-01-02"), statistics[time.Now().Format("2006-01-02")])
		if err != nil {
			log.Println(err)
		}
		newConfig, err := nascan.GetNascanConfig()
		if err != nil {
			log.Println(err)
		}
		oldConfig := configInfoAtomic.Load().(*models.ConfigNascanInfo)
		if base64.StdEncoding.EncodeToString([]byte(oldConfig.ScanList)) !=
			base64.StdEncoding.EncodeToString([]byte(newConfig.ScanList)) {
			atomic.AddInt64(&nachange, 1)
		}
		configInfoAtomic.Store(newConfig)
		time.Sleep(time.Second * 30)
	}
}

func cruise() {
	for true {

		time.Sleep(time.Hour)
	}
}

func Run() {
	fmt.Println(statisticsAtomic.Load())
	fmt.Println(configInfoAtomic.Load())
	go monitor()
	go cruise()
	for true {

		time.Sleep(time.Minute)
	}
}
