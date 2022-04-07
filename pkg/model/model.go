package model

import (
	"goblog/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	var err error

	config := mysql.New(mysql.Config{
		DSN: "luojy:secret@tcp(192.168.4.22:33306)/goblog?charset=utf8&parseTime=true&loc=Local",
	})

	DB, err = gorm.Open(config)

	logger.LogError(err)

	return DB
}

