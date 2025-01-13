package config

import (
	"fmt"
	"videohub/global"
	"videohub/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB 初始化数据库
func InitDB() {
	// 连接数据库(dns 格式)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.Mysql.Username, AppConfig.Mysql.Password, AppConfig.Mysql.Host,
		AppConfig.Mysql.Port, AppConfig.Mysql.Name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("Error connecting to mysql: %v", err)
	}
	logrus.Info("Database connected successfully")

	// 自动迁移
	db.AutoMigrate(&model.User{}, &model.Video{}, &model.Comment{}, &model.Collection{}, &model.LikeRecord{}, &model.Stats{})

	// 可以加入其他的数据库配置
	global.DB = db
}
