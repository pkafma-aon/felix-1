package util

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// JSONData 返回成功数据
func JSONData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, HTTPData{http.StatusOK, data})
}

// JSONSuccess 响应成功数据
func JSONSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, HTTPSuccess{http.StatusOK, "success"})
}

//JSONCheckBizErr 检查错误
//业务逻辑service错误的时候使用这个
func JSONCheckBizErr(c *gin.Context, err error) bool {
	flag := err != nil
	if flag {
		logrus.Error(err)
		c.JSON(http.StatusOK, HTTPError{http.StatusNoContent, err.Error()})
	}
	return flag
}

//JSONCheckMwErr gin middleware 报错
func JSONCheckMwErr(c *gin.Context, err error) bool {
	flag := err != nil
	if flag {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusOK, HTTPError{http.StatusResetContent, err.Error()})
	}
	return flag
}

//JSONCheckAuthErr 检查错误
func JSONCheckAuthErr(c *gin.Context, err error) bool {
	flag := err != nil
	if flag {
		logrus.Error(err)
		c.JSON(http.StatusOK, HTTPError{http.StatusNonAuthoritativeInfo, err.Error()})
	}
	return flag
}

//CheckBindArg 绑定参数并兼差错误
func CheckBindArg(c *gin.Context, argPtr interface{}) bool {
	err := c.ShouldBind(argPtr)
	flag := err != nil
	if flag {
		logrus.Error(err)
		c.JSON(http.StatusOK, HTTPError{http.StatusPartialContent, "参数错误:" + err.Error()})
	}
	return flag
}

//HTTPError .
type HTTPError struct {
	Code int    `json:"code" example:"203"` //203-用户认证错误 204-业务代码错误 205-中间件错误 206-参数错误
	Msg  string `json:"msg" example:"参数错误:..."`
}

//HTTPData .
type HTTPData struct {
	Code int         `json:"code" example:"200"`
	Data interface{} `json:"data"`
}

//HTTPSuccess .
type HTTPSuccess struct {
	Code int    `json:"code" example:"200"`
	Msg  string `json:"msg" example:"ok"`
}
