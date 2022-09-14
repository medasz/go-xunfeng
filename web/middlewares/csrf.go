package middlewares

import (
	"net/http"

	"go-xunfeng/config"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
)

func CSRF() gin.HandlerFunc {
	csrfMd := csrf.Protect(
		[]byte(config.Cfg.SecretKey),
		csrf.Secure(false),
		csrf.HttpOnly(true),
		csrf.ErrorHandler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte(`{"message":"CSRF token invalid"}`))
		})),
	)
	return adapter.Wrap(csrfMd)
}
