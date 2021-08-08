package model

import (
	"github.com/mitchellh/go-homedir"
	"gorm.io/driver/sqlite"
	"log"
	"math/rand"
	"os"
	"path"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

const sqliteFilePath = "/etc/felix/db.sqlite3"

func init() {
	rand.Seed(time.Now().Unix())

}

func LoadDatabase(isTest bool) *gorm.DB {
	sqliteDbPath := sqliteFilePath
	if _, err := os.Stat(sqliteFilePath); os.IsNotExist(err) {
		dir, err := homedir.Dir()
		if err != nil {
			log.Fatal("get home dir failed:", err)
		}
		sqliteDbPath = path.Join(dir, ".felix.sqlite3")
	}
	if isTest {
		sqliteDbPath = "test.db"
	}

	tempDB, err := gorm.Open(sqlite.Open(sqliteDbPath), &gorm.Config{
		SkipDefaultTransaction:                   false,
		NamingStrategy:                           nil,
		FullSaveAssociations:                     false,
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: true,
		DisableNestedTransaction:                 false,
		AllowGlobalUpdate:                        false,
		QueryFields:                              false,
		CreateBatchSize:                          0,
	})

	if err != nil {
		log.Println(err)
		return nil
	}
	err = tempDB.AutoMigrate(CfgModel{}, CfgModel{}, Machine{}, User{})
	if err != nil {
		log.Println(err)
	}
	return tempDB
}
