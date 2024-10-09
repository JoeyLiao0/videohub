package repository

import (
	"gorm.io/gorm"
)

type Comment_repository struct {
	DB *gorm.DB
}

// 工厂函数，存储单例的数据库操作对象
func NewComment_respository(db *gorm.DB) *Comment_repository {
	return &Comment_repository{DB: db}
}
