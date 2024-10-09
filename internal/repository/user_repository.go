package repository

import (
	"gorm.io/gorm"
)

type User_repository struct {
	DB *gorm.DB
}

// 工厂函数，存储单例的数据库操作对象
func NewUser_respository(db *gorm.DB) *User_repository {
	return &User_repository{DB: db}
}
