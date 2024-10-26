package repository

import (
	"fmt"

	"gorm.io/gorm"
)

type Comment struct {
	dB *gorm.DB
}

// 工厂函数，存储单例的数据库操作对象
func NewComment(db *gorm.DB) *Comment {
	return &Comment{dB: db}
}

/*
*@author:廖嘉鹏
*@create_at:2024/10/17
 */
// 测试，这是一个样板
func (cr *Comment) Test() error {
	fmt.Println("Comment_repository.Test()调用正常")
	return nil
}