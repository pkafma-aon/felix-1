package web

import (
	"fmt"
	"github.com/bytegang/felix/frontendbuild"
	"github.com/bytegang/felix/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net"
)

type Web struct {
	db  *gorm.DB
	cfg *model.CfgSchema
	lis net.Listener
}

const httpErrCode = 444

func RunWebUI(db *gorm.DB, cfg *model.CfgSchema, lis net.Listener) error {
	web := Web{
		db:  db,
		cfg: cfg,
		lis: lis,
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.MaxMultipartMemory = 32 << 20

	frontendbuild.MwServeFrontendFiles(r)
	api := r.Group("api")

	api.GET("config", web.apiConfigList())
	api.PATCH("config", web.apiConfigUpdate())

	api.GET("machine", web.apiMachineList())
	api.POST("machine", web.apiMachineCreate())
	api.PATCH("machine", web.apiMachineUpdate())
	api.DELETE("machine", web.apiMachineDelete())

	api.GET("user", web.apiUserList())
	api.POST("user", web.apiUserCreate())
	api.PATCH("user", web.apiUserUpdate())
	api.DELETE("user", web.apiUserDelete())
	api.POST("user-github-public-key-sync", web.apiUserGithubPublicKeySync())

	fmt.Println("使用 Control+C 关闭web服务")

	if err := r.RunListener(lis); err != nil {
		return err
	}
	return nil
}
