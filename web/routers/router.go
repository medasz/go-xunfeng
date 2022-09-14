package routers

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"go-xunfeng/config"
	"go-xunfeng/web/middlewares"
	"go-xunfeng/web/views"
)

func InitRouters() *gin.Engine {
	app := gin.New()
	store := sessions.NewCookieStore([]byte(config.Cfg.SecretKey))
	store.Options(sessions.Options{
		MaxAge: 6 * 60 * 60,
	})
	app.Use(
		gin.Logger(),
		gin.Recovery(),
		middlewares.CSRF(),
		sessions.Sessions("login", store),
	)

	app.GET("404", views.NotFound)
	app.GET("500", views.Error)

	app.GET("", views.GetLogin)
	app.GET("login", views.GetLogin)
	app.POST("login", views.PostLogin)
	logout := app.Group("logout")
	logout.Use(middlewares.JWT())
	{
		logout.GET("", views.LoginOut)
	}

	search := app.Group("filter")
	search.Use(middlewares.JWT())
	{
		search.GET("", views.Search)
		search.GET("search", views.Main)
		search.GET("searchxls", views.SearchResultXls)
	}

	configG := app.Group("config")
	configG.Use(middlewares.JWT())
	{
		configG.GET("", views.Config)
		configG.POST("update", views.UpdateConfig)
		configG.GET("update", views.UpdateConfig)
	}

	task := app.Group("task")
	task.Use(middlewares.JWT())
	{
		task.GET("", views.Task)
		task.POST("", views.AddTask)
		task.POST("deleteall", views.DeleteAll)
		task.POST("deletetask", views.DeleteTask)
		task.GET("deletetask", views.DeleteTask)
		task.GET("downloadxls", views.DownloadXls)
		task.POST("downloadxls", views.DownloadXls)
	}
	plugin := app.Group("plugin")
	plugin.Use(middlewares.JWT())
	{
		plugin.GET("", views.GetPlugin)
		plugin.POST("", views.GetPlugin)
		plugin.GET("list", views.Plugin)
	}

	analysis := app.Group("analysis")
	analysis.Use(middlewares.JWT())
	{
		analysis.GET("", views.Analysis)
	}

	//app.GET("loginout", views.LoginOut)
	//
	//app.GET("installplugin", views.InstallPlugin)
	//app.GET("checkupdate", views.CheckUpdate)
	//app.GET("pullupdate", views.PullUpdate)
	//app.GET("deleteplugin", views.GetDeletePlugin)
	//app.POST("deleteplugin", views.PostDeletePlugin)
	//app.GET("addplugin", views.GetAddPlugin)
	//app.POST("addplugin", views.PostAddPlugin)

	return app
}
