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

func Plugin(ctx *gin.Context) {
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
	data, err := db.GetPlugin(selector)
	if err != nil {
		log.Println(err)
		ctx.Redirect(http.StatusMovedPermanently, "/500")
		return
	}
	ctx.JSON(http.StatusOK, data)
}
