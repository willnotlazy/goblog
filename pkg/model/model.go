package model

import (
	"goblog/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	var err error

	config := mysql.New(mysql.Config{
		DSN: "luojy:secret@tcp(192.168.4.22:33306)/goblog?charset=utf8mb4&parseTime=true&loc=Local",
	})

	DB, err = gorm.Open(config, &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})

	logger.LogError(err)

	return DB
}
