package web

import (
	"net/http"

	"go-xunfeng/config"
	"go-xunfeng/pkg/htmlRender"
	"go-xunfeng/web/routers"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(config.Cfg.Mode)
}

func Web() {
	app := routers.InitRouters()
	server := http.Server{
		Addr:    config.Cfg.Server.HostPort(),
		Handler: app,
	}

	// 加载模板路径
	app.HTMLRender = htmlRender.CreateMyRender()

	// 加载静态文件，注意第一个路径参数映射第二个目录参数，所以第一个参数可以随意，但是在html中引入时需要与其保持一致
	app.Static("/static", "web/statics")

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
