package service

import (
	"videohub/internal/repository"
)

type UserList struct {
	//用户列表服务，只用到user表操作，所以只用注入user_repostiory
	userRepo *repository.User
}

// 工厂函数，返回单例的服务层操作对象
func NewUserList(ur *repository.User) *UserList {
	return &(UserList{userRepo: ur})
}
