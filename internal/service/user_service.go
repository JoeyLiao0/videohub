package service

import (
	"videohub/internal/repository"
)

type User_service struct {
	//用户服务，用到user表、collection表、video表操作
	user_repository       *repository.User_repository
	collection_repository *repository.Collection_repository
	video_repository      *repository.Video_repository
}

// 工厂函数，返回单例的服务层操作对象
func NewUser_service(ur *repository.User_repository, cr *repository.Collection_repository, vr *repository.Video_repository) *User_service {
	return &(User_service{user_repository: ur, collection_repository: cr, video_repository: vr})
}
