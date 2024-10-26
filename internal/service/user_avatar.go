package service

import (
	"fmt"
	"videohub/internal/repository"
)

type User_avatar struct {
	//头像服务，只用到user表操作，所以只用注入user_repostiory
	userRepo *repository.User
}

// 工厂函数，返回单例的服务层操作对象
func NewUser_avatar(ur *repository.User) *User_avatar {
	return &(User_avatar{userRepo: ur})
}

/*
*@author:廖嘉鹏
*@create_at:2024/10/17
 */
// 测试，这是一个样板
func (uas *User_avatar) Test() error {
	fmt.Println("User_avatar_service.Test()调用正常")
	uas.userRepo.Test()
	return nil
}

//服务函数追加在下面
