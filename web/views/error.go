package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFound(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "404", nil)
}
func Error(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "500", nil)
}
