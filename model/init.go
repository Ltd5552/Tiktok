package model

import (
	"Tiktok/config"
	"Tiktok/pkg/log"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func InitDB() {

	DB, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DatabaseSetting.User,
		config.DatabaseSetting.Password,
		config.DatabaseSetting.Host,
		config.DatabaseSetting.Port,
		config.DatabaseSetting.Name)))
	if err != nil {
		log.Fatal("fatal error connecting to the database: ", log.String("err", err.Error()))
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("fatal error setting database: ", zap.Error(err))
	}

	sqlDB.SetMaxIdleConns(config.DatabaseSetting.MaxConn)
	sqlDB.SetMaxOpenConns(config.DatabaseSetting.MaxOpen)

}
