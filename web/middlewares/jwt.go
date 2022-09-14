package middlewares

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s := sessions.Default(ctx)
		if s.Get("login") != "loginsuccess" {
			ctx.Redirect(http.StatusMovedPermanently, "/login")
			return
		}
		ctx.Next()
	}
}
