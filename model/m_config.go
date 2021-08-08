package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
)

type CfgSchema struct {
	Secret              string `desc:"APP secret 16/24/32长度" def:"uaz6225fFMTIwScSjRB842ujvQdlxyMS"`
	Version             string `desc:"版本号" def:"v0.0.1"`
	QiniuAk             string `desc:"七牛云AK" def:"xxx"`
	QiniuSk             string `desc:"七牛云SK" def:"sk"`
	QiniuBucket         string `desc:"七牛云Bucket" def:"mojocn"`
	QiniuBucketEndPoint string `desc:"七牛云Bucket endpoint" def:"https://s3.mojotv.cn"`
	TencentCloudAk      string `desc:"腾讯云AK" def:"-"`
	TencentCloudSk      string `desc:"腾讯云SK" def:"-"`
	AddrRpc             string `desc:"私有启动RPC服务地址,提供给sshd进行通讯和调用" def:"127.0.0.1:8099"`
	AddrWebSshd         string `desc:"私有启动sshd网页服务地址" def:"127.0.0.1:8025"`
	AddrWebUi           string `desc:"私有启动felix本地web管理地址" def:"127.0.0.1:8014"`
	AddrGuacamole       string `desc:"第三方guacamole gaucd服务启动的地址" def:"10.13.84.219:4822"`
}

type CfgModel struct {
	Name  string `gorm:"primaryKey" json:"name"` //配置名称
	Value string `json:"value"`                  // 配置值
	Def   string `json:"def"`                    // 配置默认值
	Desc  string `json:"desc"`                   // 配置说明
}

func CfgList(db *gorm.DB) (list []CfgModel, err error) {
	err = db.Find(&list).Error
	return
}
func CfgUpdate(db *gorm.DB, arg *CfgModel) error {
	return db.Model(new(CfgModel)).Where("name = ?", arg.Name).Updates(arg).Error
}

func CfgSyncDatabase(db *gorm.DB, res *CfgSchema) error {
	e := reflect.ValueOf(res).Elem()
	var allNames []string
	for i := 0; i < e.NumField(); i++ {
		ins := new(CfgModel)
		name := e.Type().Field(i).Name
		allNames = append(allNames, name)
		def := e.Type().Field(i).Tag.Get("def")   //get def tag of default value
		desc := e.Type().Field(i).Tag.Get("desc") //get def tag of default value
		value := e.Field(i).String()

		err := db.First(ins, "name = ?", name).Error
		if value == "" { //结构体值为空
			value = ins.Value
			if value == "" {
				value = def
			}
			e.Field(i).SetString(value)
		}

		//更新def desc字段
		ins.Def = def
		ins.Name = name
		ins.Desc = desc
		ins.Value = value

		err = db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},                            // key colume
			DoUpdates: clause.AssignmentColumns([]string{"value", "desc", "def"}), // column needed to be updated
		}).Create(ins).Error
		if err != nil {
			return err
		}
	}

	//remove redundant keys

	err := db.Delete(new(CfgModel), `name NOT IN ?`, allNames).Error
	if err != nil {
		return err
	}
	return nil
}
