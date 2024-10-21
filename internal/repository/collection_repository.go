package repository

import (
	"fmt"

	"gorm.io/gorm"
)

type Collection_repository struct {
	dB *gorm.DB
}

// 工厂函数，存储单例的数据库操作对象
func NewCollection_repository(db *gorm.DB) *Collection_repository {
	return &Collection_repository{dB: db}
}

/*
*@author:廖嘉鹏
*@create_at:2024/10/17
 */
// 测试，这是一个样板
func (cr *Collection_repository) Test() error {
	fmt.Println("Collection_repository.Test()调用正常")
	return nil
}
