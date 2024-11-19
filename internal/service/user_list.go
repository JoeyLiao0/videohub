package service

import (
	"videohub/internal/model"
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

// GetAllUsers 获取所有用户信息
func (uls *UserList) GetAllUsers() ([]model.User, error) {
	// 调用 user_repository 获取所有用户的信息
	users, err := uls.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}
