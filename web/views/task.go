package views

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	"go.mongodb.org/mongo-driver/bson"

	"go-xunfeng/db"
	"go-xunfeng/models"
	"go-xunfeng/pkg/query"
	"go-xunfeng/web/params"
)

func Task(ctx *gin.Context) {
	param := new(params.TaskList)
	if err := ctx.ShouldBind(&param); err != nil {
		log.Println(err)
		ctx.Redirect(http.StatusOK, "/500")
		return
	}
	total, err := db.CountAllTask(bson.M{})
	if err != nil {
		log.Println(err)
		ctx.Redirect(http.StatusMovedPermanently, "/500")
		return
	}
	items, err := db.GetTask(bson.M{}, param.Page, 60)
	if err != nil {
		log.Println(err)
		ctx.Redirect(http.StatusMovedPermanently, "/500")
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
	ctx.HTML(http.StatusOK, "task", bson.M{
		"total":     total,
		"item":      items,
		"page_list": pageList,
		"csrf":      csrf.Token(ctx.Request),
	})
}

func AddTask(ctx *gin.Context) {
	param := new(params.Task)
	if err := ctx.ShouldBind(&param); err != nil {
		log.Println(err)
		ctx.Redirect(http.StatusOK, "/500")
		return
	}

	result := "fail"
	if param.Plugin != "" {
		targets := make([][]interface{}, 0)
		q, err := query.QueryLogic(param.Condition)
		if err != nil {
			log.Println(err)
			ctx.String(http.StatusOK, result)
			return
		}
		qString, err := json.Marshal(q)
		if err != nil {
			log.Println(err)
			ctx.String(http.StatusOK, result)
			return
		}
		if param.ResultCheck {
			targets, err = db.GetInfoAllIpPort(q)
			if err != nil {
				log.Println(err)
				ctx.String(http.StatusOK, result)
				return
			}
		} else {
			ids := strings.Split(param.Ids, ",")
			for _, i := range ids {
				res := strings.Split(i, ":")
				if len(res) == 2 {
					targets = append(targets, []interface{}{res[0], res[1]})
				}
			}
		}

		tmpResult := true
		pluginList := strings.Split(param.Plugin, ",")
		for _, plugin := range pluginList {
			task := models.InTask{
				Title:     param.Title,
				Status:    0,
				IsUpdate:  param.IsUpdate,
				Time:      time.Now(),
				Query:     string(qString),
				Plan:      0,
				Plugin:    plugin,
				Condition: param.Condition,
				Target:    targets,
			}
			res, err := db.CreateTask(task)
			if err != nil {
				log.Println(err)
				ctx.String(http.StatusOK, result)
				return
			}
			if res == nil {
				tmpResult = false
			}
		}
		if tmpResult {
			result = "success"
		}
	}
	ctx.String(http.StatusOK, result)
}

func DeleteAll(ctx *gin.Context) {
	err := db.TaskDeleteAll()
	if err != nil {
		log.Println(err)
		ctx.String(http.StatusOK, "fail")
		return
	}
	ctx.String(http.StatusOK, "success")
}

func DeleteTask(ctx *gin.Context) {
	param := new(params.TaskId)
	if err := ctx.ShouldBind(&param); err != nil {
		log.Println(err)
		ctx.Redirect(http.StatusOK, "/500")
		return
	}
	if param.Oid != "" {
		res, err := db.DeleteTaskById(param.Oid)
		if err != nil {
			log.Println(err)
			ctx.String(http.StatusOK, "fail")
			return
		}
		if res.DeletedCount > 0 {
			delRes, err := db.DeleteResultByTaskId(param.Oid)
			if err != nil {
				log.Println(err)
				ctx.String(http.StatusOK, "fail")
				return
			}
			if delRes != nil {
				ctx.String(http.StatusOK, "success")
				return
			}
		}
	}
	ctx.String(http.StatusOK, "fail")
}

func DownloadXls(ctx *gin.Context) {
	param := new(params.TaskDownloadXls)
	if err := ctx.ShouldBind(&param); err != nil {
		log.Println(err)
		ctx.Redirect(http.StatusOK, "/500")
		return
	}
	fmt.Println(param)
	if param.TaskId == "" {

	} else {
		//下载综合报表
		//task, err := db.GetTaskAll()
		//if err != nil {
		//	log.Println(err)
		//	ctx.Redirect(http.StatusOK, "/500")
		//	return
		//}
	}
}
