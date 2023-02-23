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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DatabaseSetting.User,
		config.DatabaseSetting.Password,
		config.DatabaseSetting.Host,
		config.DatabaseSetting.Port,
		config.DatabaseSetting.Name)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("fatal error connecting to the database: ", zap.Error(err))
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("fatal error setting database: ", zap.Error(err))
	}

	log.Info("Gorm init successfully")
	sqlDB.SetMaxIdleConns(config.DatabaseSetting.MaxConn)
	sqlDB.SetMaxOpenConns(config.DatabaseSetting.MaxOpen)
}
