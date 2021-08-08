package web

import (
	"github.com/bytegang/felix/model"
	"github.com/gin-gonic/gin"
)

func (w *Web) apiConfigList() gin.HandlerFunc {
	return func(c *gin.Context) {
		list, err := model.CfgList(w.db)
		if err != nil {
			c.JSON(httpErrCode, err)
			return
		}
		c.JSON(200, list)
	}
}

func (w *Web) apiConfigUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		arg := new(model.CfgModel)
		err := c.BindJSON(arg)
		if err != nil {
			c.JSON(httpErrCode, err)
			return
		}
		err = model.CfgUpdate(w.db, arg)
		if err != nil {
			c.JSON(httpErrCode, err)
			return
		}
		c.JSON(200, "ok")
	}

}
