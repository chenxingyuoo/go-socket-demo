package router

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func RouterInit (r *gin.Engine) {
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello World!")
	})
}