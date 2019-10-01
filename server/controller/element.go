package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-socket-demo/socket"
	"go-socket-demo/common"
	"go-socket-demo/server"
)

type SaveReqBody struct{
	Token string  `json:"token"`
	ProjectId int `json:"projectId"`
	ElementId int `json:"elementId"`
}

type ResBody struct {
	Url string `json:"url"`
	SaveReqBody
}

// 请求截图服务
func requestScreenshot(reqInfo SaveReqBody)  {
	projectId := strconv.Itoa(reqInfo.ProjectId)
	elementId := strconv.Itoa(reqInfo.ElementId)
	previewUrl := fmt.Sprintf("http://localhost:8082/?projectId=%s&elementId=%s", projectId, elementId)
	callbackUrl := "http://localhost:8000/element/screenshotCallback"

	param := make(map[string]interface{})
	param["token"] = reqInfo.Token
	param["projectId"] = reqInfo.ProjectId
	param["elementId"] = reqInfo.ElementId
	param["previewUrl"] = previewUrl
	param["callbackUrl"] = callbackUrl

	server.Screenshot(param)
}

// 保存元素
func Save(c *gin.Context)  {
  var reqInfo SaveReqBody
	if err := c.ShouldBindJSON(&reqInfo); err != nil  {
		 c.JSON(http.StatusBadRequest, common.ResponseError(err.Error()))
		 return
	}

	token := c.GetHeader("token")
	reqInfo.Token = token

	fmt.Println("reqInfo",reqInfo)

	go requestScreenshot(reqInfo)

	c.JSON(200, common.ResponseSuccess(nil))
}

// 截图成功回调
func ScreenshotCallback(c *gin.Context)  {
	var resBody ResBody
	if err := c.ShouldBindJSON(&resBody); err != nil  {
		 c.JSON(http.StatusBadRequest, common.ResponseError(err.Error()))
		 return
	}

	// fmt.Println("resBody", resBody)

	param := make(map[string]interface{})
	data := make(map[string]interface{})
	data["url"] = resBody.Url
	data["elementId"] = resBody.ElementId
	data["projectId"] = resBody.ProjectId
	
	param["type"] = "thumbnail"
	param["data"] = data

	bytesData, err := json.Marshal(param)
	if err != nil {
		fmt.Println(err.Error() )
		return
	}
	
	// 向客户端发送消息
	wsConn := socket.GetServer(resBody.Token)
	wsConn.WsWrite(websocket.TextMessage, bytesData)

	c.JSON(200, common.ResponseSuccess(nil))
}