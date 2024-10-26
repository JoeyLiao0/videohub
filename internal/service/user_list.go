package service

import (
	"fmt"
	"videohub/internal/repository"
)

type User_list struct {
	//用户列表服务，只用到user表操作，所以只用注入user_repostiory
	userRepo *repository.User
}

// 工厂函数，返回单例的服务层操作对象
func NewUser_list(ur *repository.User) *User_list {
	return &(User_list{userRepo: ur})
}

/*
*@author:廖嘉鹏
*@create_at:2024/10/17
 */
// 测试，这是一个样板
func (uls *User_list) Test() error {
	fmt.Println("User_list_service.Test()调用正常")
	uls.userRepo.Test()
	return nil
}

//服务函数追加在下面
