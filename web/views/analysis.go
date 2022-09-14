package views

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"go-xunfeng/db"
)

func Analysis(ctx *gin.Context) {
	ip, err := db.GetInfoIpCount()
	if err != nil {
		log.Println(err)
	}
	record, err := db.CountAllInfo(bson.M{})
	if err != nil {
		log.Println(err)
	}
	task, err := db.CountAllTask(bson.M{})
	if err != nil {
		log.Println(err)
	}
	vul, err := db.GetVulCount()
	if err != nil {
		log.Println(err)
	}
	plugin, err := db.CountPluginAll(bson.M{})
	if err != nil {
		log.Println(err)
	}
	vulType, err := db.GetVulGroupType()
	if err != nil {
		log.Println(err)
	}
	trend, err := db.GetStatisticsLimit()
	if err != nil {
		log.Println(err)
	}
	vulBeat, err := db.GetHeartbeat("load")
	if err != nil {
		log.Println(err)
	}
	scanBeat, err := db.GetHeartbeat("heartbeat")
	if err != nil {
		log.Println(err)
	}
	taskPercent := vulBeat.Value * 100
	taskAliveSecond := time.Now().Sub(vulBeat.UpTime).Seconds()
	scanAliveSecond := time.Now().Sub(scanBeat.UpTime).Seconds()
	var taskAlive, scanAlive bool
	if taskAliveSecond < 120 {
		taskAlive = true
	}
	if scanAliveSecond < 120 {
		scanAlive = true
	}
	serverType, err := db.GetServerType()
	if err != nil {
		log.Println(err)
	}
	webType, err := db.GetWebType()
	if err != nil {
		log.Println(err)
	}
	ctx.HTML(http.StatusOK, "analysis", bson.M{
		"ip":          ip,
		"record":      record,
		"task":        task,
		"vul":         vul,
		"plugin":      plugin,
		"vultype":     vulType,
		"trend":       trend,
		"taskpercent": taskPercent,
		"taskalive":   taskAlive,
		"scanalive":   scanAlive,
		"server_type": serverType,
		"web_type":    webType,
	})
}
