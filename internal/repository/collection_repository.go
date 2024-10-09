package repository

import (
	"gorm.io/gorm"
)

type Collection_repository struct {
	DB *gorm.DB
}

// 工厂函数，存储单例的数据库操作对象
func NewCollection_respository(db *gorm.DB) *Collection_repository {
	return &Collection_repository{DB: db}
}
