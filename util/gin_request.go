package util

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

//QueryInt 解析query中的参数到int
func QueryInt(c *gin.Context, key string) (int, error) {
	is := c.Query(key)
	number, err := strconv.ParseInt(is, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(number), nil
}

//ParamUint 解析path-param中的参数到uint
func ParamUint(c *gin.Context, key string) (uint, error) {
	is := c.Param(key)
	number, err := strconv.ParseUint(is, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(number), nil
}
