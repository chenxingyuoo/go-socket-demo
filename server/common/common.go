package common

import (
	"github.com/gin-gonic/gin"
)

// 成功
func ResponseSuccess(data interface{}) (res *gin.H) {
	res = &gin.H{
		"data": data,
		"code": 200,
		"message": "success",
	}

	return 
}

// 失败
func ResponseFail(data interface{}) (res *gin.H) {
	res = &gin.H{
		"data": data,
		"code": 500,
		"message": "success",
	}

	return 
}

// 错误
func ResponseError(err interface{}) (res *gin.H) {
	res = &gin.H{
		"error": err,
	}

	return 
}