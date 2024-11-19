package service

import (
	"videohub/internal/repository"
)

type VideoList struct {
	//视频列表服务，用到user表、video表操作
	userRepo  *repository.User
	videoRepo *repository.Video
}

// 工厂函数，返回单例的服务层操作对象
func NewVideoList(ur *repository.User, vr *repository.Video) *VideoList {
	return &(VideoList{userRepo: ur, videoRepo: vr})
}

//服务函数追加在下面
