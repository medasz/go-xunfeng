package vulscan

import (
	"fmt"
	"time"

	"github.com/medasz/kunpeng"
	"go.mongodb.org/mongo-driver/bson"

	"go-xunfeng/db"
	"go-xunfeng/models"
)

func init() {
	pluginCount, err := db.CountPluginAll(bson.M{})
	if err != nil {
		panic(err)
	}
	if pluginCount >= 1 {
		return
	}
	levelList := []string{"紧急", "高危", "中危", "低危", "提示"}
	plugins := kunpeng.Greeter.GetPlugins()
	now := time.Now()
	for _, plugin := range plugins {
		keyword := ""
		if _, ok := plugin["keyword"]; ok {
			keyword = plugin["keyword"].(string)
		}
		pluginInfo := models.Plugin{
			Id:       plugin["references"].(map[string]interface{})["kpid"].(string),
			Count:    0,
			AddTime:  now,
			Info:     fmt.Sprintf("%s %s", plugin["remarks"], plugin["references"].(map[string]interface{})["cve"].(string)),
			Name:     fmt.Sprintf("%s%s", "Kunpeng -", plugin["references"].(map[string]interface{})["name"]),
			Keyword:  keyword,
			Level:    levelList[int(plugin["level"].(float64))],
			Url:      plugin["references"].(map[string]interface{})["url"].(string),
			Author:   plugin["author"].(string),
			Filename: plugin["references"].(map[string]interface{})["kpid"].(string),
			Source:   1,
			Type:     plugin["type"].(string),
		}
		_, err := db.CreatePlugin(pluginInfo)
		if err != nil {
			panic(err)
		}
	}
}

func Run() {
	println("开始...", len(kunpeng.Greeter.GetPlugins()))
}
