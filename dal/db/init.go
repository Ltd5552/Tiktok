package db

import (
	"Tiktok/pkg/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func init() {
	var err error
	dsn := ""
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil{
		log.Error("sql open error")
	}
}