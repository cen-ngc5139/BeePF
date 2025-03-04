package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success   bool        `json:"success"`
	ErrorCode int         `json:"errorCode"`
	ErrorMsg  string      `json:"errorMsg"`
	Data      interface{} `json:"data,omitempty"`
}

func ResponseOk(c *gin.Context, data interface{}) {
	resp := &Response{
		Success:   true,
		ErrorCode: http.StatusOK,
		ErrorMsg:  "",
		Data:      data,
	}
	c.JSON(http.StatusOK, resp)
}

func ResponseErrorCode(c *gin.Context, errCode int, errMsg string) {
	resp := &Response{
		Success:   false,
		ErrorCode: errCode,
		ErrorMsg:  errMsg,
		Data:      nil,
	}
	c.JSON(http.StatusOK, resp)
}

func ResponseErrorCodeWithHttpCode(c *gin.Context, errCode int, errMsg string, httpCode int) {
	resp := &Response{
		Success:   false,
		ErrorCode: errCode,
		ErrorMsg:  errMsg,
		Data:      nil,
	}
	c.JSON(httpCode, resp)
	c.Abort()
}
