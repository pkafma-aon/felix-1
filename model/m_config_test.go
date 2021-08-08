package model

import (
	"testing"
)

func TestCfgSyncDatabase(t *testing.T) {
	db := LoadDatabase(true)
	cfg := new(CfgSchema)
	vV := "v.2.2.4"
	cfg.Version = vV
	CfgSyncDatabase(db, cfg)

	some := new(CfgModel)
	db.First(some, "name = ?", "Version")
	if some.Value != vV || cfg.Version != vV {
		t.Error("err")
	}

	some = new(CfgModel)
	db.First(some, "name = ?", "QiniuAk")
	vV = "xxx"
	if some.Value != vV || cfg.QiniuAk != vV || some.Desc != "七牛云AK" {
		t.Error("err")
	}
}
