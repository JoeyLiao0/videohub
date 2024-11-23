package config

import (
	"fmt"
	"log"
	"videohub/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	// 或者直接将配置文件的信息改为 dsn
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", 
					AppConfig.Mysql.Username, AppConfig.Mysql.Password, AppConfig.Mysql.Host, 
					AppConfig.Mysql.Port, AppConfig.Mysql.Name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error connecting to mysql: %v", err)
	}

	log.Println("Connected to mysql")

	db.AutoMigrate(&model.User{}, &model.Video{}, &model.VideoChunk{})

	// todo: 数据库的一些其他设置
	return db
}
