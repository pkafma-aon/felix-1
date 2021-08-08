package web

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytegang/felix/model"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

func (w Web) apiUserList() gin.HandlerFunc {
	return func(c *gin.Context) {
		q := c.Query("q")
		var list []model.User
		tx := w.db.Model(new(model.User))
		if q != "" {
			tx.Where("`name` LIKE ? OR `email` LIKE ?", "%"+q+"%", "%"+q+"%")
		}
		err := w.db.Find(&list).Error
		if err != nil {
			c.JSON(httpErrCode, err)
			return
		}
		c.JSON(200, list)
	}
}

func (w Web) apiUserCreate() gin.HandlerFunc {
	return func(c *gin.Context) {
		arg := new(model.User)
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
func (w Web) apiUserUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		arg := new(model.User)
		err := c.BindJSON(arg)
		if err != nil {
			c.JSON(httpErrCode, err)
			return
		}
		err = w.db.Model(new(model.User)).Where("id = ?", arg.Id).Updates(arg).Error
		if err != nil {
			c.JSON(httpErrCode, err)
			return
		}
		c.JSON(200, "ok")
	}
}
func (w Web) apiUserDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		arg := new(model.User)
		err := c.BindJSON(arg)
		if err != nil {
			c.JSON(httpErrCode, err)
			return
		}

		err = w.db.Delete(new(model.User), arg.Id).Error
		if err != nil {
			c.JSON(httpErrCode, err)
			return
		}
		c.JSON(200, "ok")
	}
}

//apiUserGithubPublicKeySync 同步github public key
func (w Web) apiUserGithubPublicKeySync() gin.HandlerFunc {
	return func(c *gin.Context) {
		arg := new(model.User)
		err := c.BindJSON(arg)
		if err != nil {
			c.JSON(httpErrCode, err)
			return
		}

		err = w.db.First(arg, arg.Id).Error
		if err != nil {
			c.JSON(httpErrCode, err)
			return
		}
		if arg.GithubAccount == "" {
			c.JSON(httpErrCode, "your github account is empty")
			return
		}
		keys, err := fetchGithubPublicKeys(arg.GithubAccount)
		if err != nil {
			c.JSON(httpErrCode, err)
			return
		}
		arg.PublicKey = keys
		err = w.db.Model(arg).Select("PublicKey").Updates(arg).Error
		if err != nil {
			c.JSON(httpErrCode, err)
			return
		}
		c.JSON(200, "ok")
	}
}

func fetchGithubPublicKeys(githubUser string) (string, error) {
	keyURL := fmt.Sprintf("https://github.com/%s.keys", githubUser)
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*15)
	defer cancelFunc()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, keyURL, nil)
	if err != nil {
		return "", err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", errors.New("invalid response from github")
	}
	authorizedKeysBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("reading body:%v", err)
	}
	return string(authorizedKeysBytes), nil
}
