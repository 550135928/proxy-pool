package databases

import (
	"log"
	"proxy-pool/config"
	"proxy-pool/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// ORM 结构体
type ORM struct {
	DB *gorm.DB
}

// New 数据库使用sqlite
func New(config *config.Config) *ORM {

	var err error
	var db *gorm.DB

	db, err = gorm.Open("sqlite3", config.Sqlite.DBFilePath)
	// TODO: to config file
	db.LogMode(true)
	if err != nil {
		log.Fatalf("mysql connect error %v", err)
	}
	// 创建表结构
	if err := db.AutoMigrate(&model.Proxy{}).Error; err != nil {
		log.Fatalf("mysql create table err:%#v", err)
	}

	return &ORM{DB: db}
}

// Close close
func (o *ORM) Close() {
	o.DB.Close()
}
