package web

import (
	"fmt"
	"github.com/bytegang/felix/model"
	"github.com/gin-gonic/gin"
	"time"
)

func (w Web) apiMachineList() gin.HandlerFunc {
	return func(c *gin.Context) {
		q := c.Query("q")
		var list []model.Machine
		tx := w.db.Model(new(model.Machine))
		if q != "" {
			tx.Where("name LIKE ? OR host LIKE ?", "%"+q+"%", "%"+q+"%")
		}
		err := w.db.Find(&list).Error
		if err != nil {
			c.JSON(httpErrCode, err)
			return
		}
		//generate web ssh terminal URL
		for i, machine := range list {
			token, err := machine.GenerateToken([]byte(w.cfg.Secret), time.Hour*24)
			if err != nil {
				c.JSON(httpErrCode, err)
				return
			}
			if machine.Protocol == "rdp" {
				machine.WebSshURL = fmt.Sprintf("http://%s/#/guacamole/%s", w.cfg.AddrWebSshd, token)
			} else {
				machine.WebSshURL = fmt.Sprintf("http://%s/#/ssh/%s", w.cfg.AddrWebSshd, token)
			}
			list[i] = machine
		}
		c.JSON(200, list)
	}
}

func (w Web) apiMachineCreate() gin.HandlerFunc {
	return func(c *gin.Context) {
		arg := new(model.Machine)
		err := c.BindJSON(arg)
		if err != nil {
			c.JSON(httpErrCode, err)
			return
		}
		err = w.db.Create(arg).Error
		if err != nil {
			c.JSON(httpErrCode, err)
			return
		}
		c.JSON(200, "ok")
	}
}
func (w Web) apiMachineUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		arg := new(model.Machine)
		err := c.BindJSON(arg)
		if err != nil {
			c.JSON(httpErrCode, err)
			return
		}
		err = w.db.Model(new(model.Machine)).Where("id = ?", arg.Id).Updates(arg).Error
		if err != nil {
			c.JSON(httpErrCode, err)
			return
		}
		c.JSON(200, "ok")
	}
}
func (w Web) apiMachineDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		arg := new(model.Machine)
		err := c.BindJSON(arg)
		if err != nil {
			c.JSON(httpErrCode, err)
			return
		}

		err = w.db.Delete(new(model.Machine), arg.Id).Error
		if err != nil {
			c.JSON(httpErrCode, err)
			return
		}
		c.JSON(200, "ok")
	}
}
