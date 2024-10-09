package service

import (
	"videohub/internal/repository"
)

type User_avatar_service struct {
	//头像服务，只用到user表操作，所以只用注入user_repostiory
	user_repository *repository.User_repository
}

// 工厂函数，返回单例的服务层操作对象
func NewUser_avatar_service(ur *repository.User_repository) *User_avatar_service {
	return &(User_avatar_service{user_repository: ur})
}
