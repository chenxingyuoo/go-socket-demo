package router

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"go-socket-demo/controller"
	"go-socket-demo/socket"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ws/:token", func (c *gin.Context)  {
		token := c.Param("token")
		
		socket.InitConnection(token, c.Writer, c.Request)
	})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})

	element := r.Group("/element")
	{
		element.POST("/save", controller.Save)
		element.POST("/screenshotCallback", controller.ScreenshotCallback)
	}

	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", "")
	})
	return r
}