package service

import (
	"fmt"
	"videohub/internal/repository"
)

type User_list_service struct {
	//用户列表服务，只用到user表操作，所以只用注入user_repostiory
	user_repository *repository.User_repository
}

// 工厂函数，返回单例的服务层操作对象
func NewUser_list_service(ur *repository.User_repository) *User_list_service {
	return &(User_list_service{user_repository: ur})
}

/*
*@author:廖嘉鹏
*@create_at:2024/10/17
 */
// 测试，这是一个样板
func (uls *User_list_service) Test() error {
	fmt.Println("User_list_service.Test()调用正常")
	uls.user_repository.Test()
	return nil
}

//服务函数追加在下面
