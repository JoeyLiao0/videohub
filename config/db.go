package config

import (
	"fmt"
	"videohub/global"
	"videohub/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() {
	// 或者直接将配置文件的信息改为 dsn
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.Mysql.Username, AppConfig.Mysql.Password, AppConfig.Mysql.Host,
		AppConfig.Mysql.Port, AppConfig.Mysql.Name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		logrus.Fatalf("Error connecting to mysql: %v", err)
	}

	logrus.Info("Database connected successfully")

	db.AutoMigrate(&model.User{}, &model.Video{}, &model.Comment{})

	// todo: 数据库的一些其他设置
	global.DB = db
}
