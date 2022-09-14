package views

import (
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"

	"go-xunfeng/db"

	"github.com/gin-gonic/gin"

	"go-xunfeng/web/params"
)

func InstallPlugin(ctx *gin.Context) {

}

func CheckUpdate(ctx *gin.Context) {

}

func PullUpdate(ctx *gin.Context) {

}

func GetDeletePlugin(ctx *gin.Context) {

}

func PostDeletePlugin(ctx *gin.Context) {

}

func GetAddPlugin(ctx *gin.Context) {

}

func PostAddPlugin(ctx *gin.Context) {

}

func GetPlugin(ctx *gin.Context) {
	param := new(params.Plugin)
	if err := ctx.ShouldBind(&param); err != nil {
		log.Println(err)
		ctx.Redirect(http.StatusMovedPermanently, "/500")
		return
	}
	selector := bson.M{}
	if param.Type != "" {
		selector["type"] = param.Type
	}
	if param.Risk != "" {
		selector["level"] = param.Risk
	}
	if param.Search != "" {
		selector["name"] = bson.M{"$regex": param.Search, "$options": "i"}
	}
	data, err := db.GetPluginAll(selector)
	if err != nil {
		log.Println(err)
		ctx.Redirect(http.StatusMovedPermanently, "/500")
		return
	}
	ctx.JSON(http.StatusOK, data)
}

func Plugin(ctx *gin.Context) {
	param := new(params.PluginList)
	if err := ctx.ShouldBind(&param); err != nil {
		log.Println(err)
		ctx.Redirect(http.StatusMovedPermanently, "/500")
		return
	}
	selector := bson.M{}
	data, err := db.GetPlugin(selector, param.Page, 60)
	if err != nil {
		log.Println(err)
		ctx.Redirect(http.StatusMovedPermanently, "/500")
		return
	}
	pluginTypes, err := db.GetPluginTypes()
	if err != nil {
		log.Println(err)
		ctx.Redirect(http.StatusFound, "/500")
		return
	}
	total, err := db.CountPluginAll(bson.M{})
	if err != nil {
		log.Println(err)
		ctx.Redirect(http.StatusFound, "/500")
		return
	}
	mod := total % 60
	num := 0
	if mod == 0 {
		num = total / 60
	} else {
		num = (total / 60) + 1
	}
	pageList := make([]int, 0)
	for i := 1; i <= num; i++ {
		pageList = append(pageList, i)
	}
	result := make(map[string]interface{})
	result["count"] = total
	result["data"] = data
	result["vultype"] = pluginTypes
	result["page_list"] = pageList
	ctx.HTML(http.StatusOK, "plugin", result)
}
