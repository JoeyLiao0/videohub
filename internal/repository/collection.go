package repository

import (
	"fmt"

	"gorm.io/gorm"
)

type Collection struct {
	DB *gorm.DB
}

// 工厂函数，存储单例的数据库操作对象
func NewCollection(db *gorm.DB) *Collection {
	return &Collection{DB: db}
}

/*
*@author:廖嘉鹏
*@create_at:2024/10/17
 */
// 测试，这是一个样板
func (cr *Collection) Test() error {
	fmt.Println("Collection_repository.Test()调用正常")
	return nil
}
