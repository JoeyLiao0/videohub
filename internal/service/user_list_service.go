package service

import (
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
