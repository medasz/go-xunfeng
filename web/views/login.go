package views

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"

	"go-xunfeng/config"
	"go-xunfeng/web/params"
)

func GetLogin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login", gin.H{csrf.TemplateTag: csrf.TemplateField(ctx.Request)})
}

func PostLogin(ctx *gin.Context) {
	var param *params.Login
	err := ctx.ShouldBind(&param)
	if err != nil {
		panic(err)
	}
	if param.Account == config.Cfg.User.Account && param.Password == config.Cfg.User.Password {
		s := sessions.Default(ctx)
		s.Set("login", "loginsuccess")
		s.Save()
		ctx.Redirect(http.StatusMovedPermanently, "filter")
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "login")
	}
}

func LoginOut(ctx *gin.Context) {
	s := sessions.Default(ctx)
	s.Clear()
	s.Save()
	ctx.Redirect(http.StatusMovedPermanently, "login")
}
