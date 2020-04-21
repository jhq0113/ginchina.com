package helper

import (
	"github.com/gin-gonic/gin"
)

func SetCookie(ctx *gin.Context, key string, value string, timeout int, args ...interface{}) {
	path := "/"
	domain := ctx.Request.Host
	secure := false
	if len(args) > 0 {
		path = args[0].(string)
	}

	if len(args) > 1 {
		domain = args[1].(string)
	}

	if len(args) > 2 {
		secure = args[2].(bool)
	}

	ctx.SetCookie(key, value, timeout, path, domain, secure, true)
}
