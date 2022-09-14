package views

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	"go.mongodb.org/mongo-driver/bson"

	"go-xunfeng/db"
	"go-xunfeng/models"
	"go-xunfeng/pkg/exchange"
	"go-xunfeng/web/params"
)

func Config(ctx *gin.Context) {
	param := new(params.ConfigType)
	err := ctx.BindQuery(param)
	if err != nil {
		log.Println(err)
		ctx.Redirect(http.StatusMovedPermanently, "/500")
		return
	}
	res := make([]models.ConfigExternal, 0)
	if param.Config == "nascan" {
		data, err := db.GetNascanConfig()
		if err != nil {
			log.Println(err)
			ctx.Redirect(http.StatusMovedPermanently, "/500")
			return
		}
		res = exchange.NascanConfigMongoToTemplate(data)
	} else if param.Config == "vulscan" {
		data, err := db.GetVulscanConfig()
		if err != nil {
			log.Println(err)
			ctx.Redirect(http.StatusMovedPermanently, "/500")
			return
		}
		res = exchange.VulscanConfigMongoToTemplate(data)
	}
	ctx.HTML(http.StatusOK, "config", gin.H{"values": res, "csrf": csrf.Token(ctx.Request)})
}

func UpdateConfig(ctx *gin.Context) {
	param := new(params.UpdateConfig)
	if err := ctx.ShouldBind(&param); err != nil {
		log.Println(err)
		ctx.Redirect(http.StatusFound, "/500")
		return
	}
	nascan, err := db.GetNascanConfig()
	if err != nil {
		log.Println(err)
		ctx.String(http.StatusOK, "fail")
		return
	}

	selector := bson.M{"type": param.ConfigType}
	data := bson.M{}
	if param.Name == "Masscan" {
		masscanConfig := strings.Split(nascan.ConfigNascan.Masscan.Value, "|")
		data["$set"] = bson.M{
			fmt.Sprintf("config.%s.value", param.Name): fmt.Sprintf("%s|%s", masscanConfig[0], param.Value),
		}
	} else if param.Name == "Port_list" {
		portListConfig := strings.Split(nascan.ConfigNascan.PortList.Value, "|")
		data["$set"] = bson.M{
			fmt.Sprintf("config.%s.value", param.Name): fmt.Sprintf("%s|%s", portListConfig[0], param.Value),
		}
	} else if param.Name == "Port_list_Flag" {
		portListConfig := strings.SplitN(nascan.ConfigNascan.PortList.Value, "|", 2)
		data["$set"] = bson.M{
			fmt.Sprintf("config.%s.value", "Port_list"): fmt.Sprintf("%s|%s", param.Value, portListConfig[1]),
		}
	} else if param.Name == "Masscan_Flag" {
		masscanConfig := strings.SplitN(nascan.ConfigNascan.Masscan.Value, "|", 2)
		data["$set"] = bson.M{
			fmt.Sprintf("config.%s.value", "Masscan"): fmt.Sprintf("%s|%s", param.Value, masscanConfig[1]),
		}
	} else {
		data["$set"] = bson.M{
			fmt.Sprintf("config.%s.value", param.Name): param.Value,
		}
	}
	err = db.UpdateConfig(selector, data)
	if err != nil {
		log.Println(err)
		ctx.String(http.StatusOK, "fail")
		return
	}
	ctx.String(http.StatusOK, "success")
}
