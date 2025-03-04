package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var EmptyData = &map[string]interface{}{}

func HandleResult(c *gin.Context, data interface{}) {
	jsonResult(c, data)
}

func HandleError(c *gin.Context, err error) bool {
	if err != nil {
		jsonError(c, err.Error())
		return true
	}
	return false
}

func HandleErrorCode(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(200, Response{
		Success:   false,
		ErrorCode: code,
		ErrorMsg:  msg,
		Data:      "",
	})
}

func jsonError(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(200, Response{
		Success:   false,
		ErrorCode: 500,
		ErrorMsg:  msg,
		Data:      "",
	})
}

func jsonResult(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Success:   true,
		ErrorCode: 0,
		ErrorMsg:  "",
		Data:      data,
	})
}

func ResponseErr(c *gin.Context, err error) {
	jsonError(c, err.Error())
}

func ResponseProbeError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, Response{
		Success:   false,
		ErrorCode: 0,
		ErrorMsg:  err.Error(),
		Data:      "",
	})
}

func GetParamIntItem(key string, c *gin.Context) (num int) {
	var err error
	str := c.Param(key)
	num, err = strconv.Atoi(str)
	if err != nil {
		return
	}

	return
}

func GetQueryIntItem(key, def string, c *gin.Context) (num int) {
	var err error
	str := c.DefaultQuery(key, def)
	num, err = strconv.Atoi(str)
	if err != nil {
		return
	}

	return
}

func GetPageInfo(c *gin.Context) (pageSize int, pageNum int) {
	return GetQueryIntItem("pageSize", "0", c), GetQueryIntItem("pageNum", "1", c)
}

func HandleData(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Success:   true,
		ErrorCode: 0,
		ErrorMsg:  "",
		Data:      data,
	})
}
