package repository

import (
	"gorm.io/gorm"
)

type Video_repository struct {
	DB *gorm.DB
}

// 工厂函数，存储单例的数据库操作对象
func NewVideo_respository(db *gorm.DB) *Video_repository {
	return &Video_repository{DB: db}
}
