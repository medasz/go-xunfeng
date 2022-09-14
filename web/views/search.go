package views

import (
	"log"
	"net/http"

	"go-xunfeng/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	"go.mongodb.org/mongo-driver/bson"

	"go-xunfeng/db"
	"go-xunfeng/pkg/query"
	"go-xunfeng/web/params"
)

func Search(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "search", nil)
}

func Main(ctx *gin.Context) {
	param := new(params.Search)
	if err := ctx.ShouldBind(&param); err != nil {
		log.Println(err)
		ctx.Redirect(http.StatusMovedPermanently, "/500")
		return
	}

	if param.Q == "" {
		ctx.HTML(http.StatusOK, "main", bson.M{
			"item":        []models.Info{},
			"plugin":      []models.Plugin{},
			"itemcount":   0,
			"plugin_type": []string{},
			"page_list":   []int{},
			"csrf":        csrf.Token(ctx.Request),
		})
		return
	}

	plugins, err := db.GetPluginAll(bson.M{})
	if err != nil {
		log.Println(err)
		ctx.Redirect(http.StatusFound, "/500")
		return
	}
	pluginTypes, err := db.GetPluginTypes()
	if err != nil {
		log.Println(err)
		ctx.Redirect(http.StatusFound, "/500")
		return
	}
	selector, err := query.QueryLogic(param.Q)
	if err != nil {
		log.Println(err)
		ctx.Redirect(http.StatusFound, "/500")
		return
	}
	infos, err := db.GetInfo(selector, param.Page, 60)
	if err != nil {
		log.Println(err)
		ctx.Redirect(http.StatusFound, "/500")
		return
	}
	total, err := db.CountAllInfo(selector)
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
	ctx.HTML(http.StatusOK, "main", bson.M{
		"item":        infos,
		"plugin":      plugins,
		"itemcount":   total,
		"plugin_type": pluginTypes,
		"page_list":   pageList,
		"csrf":        csrf.Token(ctx.Request),
	})
}

// SearchResultXls 搜索结果报表下载接口
func SearchResultXls(ctx *gin.Context) {

}
