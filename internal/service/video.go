package service

import (
	"videohub/internal/repository"
)

type Video struct {
	//视频服务，用到user表、video表、comment表操作
	userRepo    *repository.User
	videoRepo   *repository.Video
	commentRepo *repository.Comment
}

// 工厂函数，返回单例的服务层操作对象
func NewVideo(ur *repository.User, vr *repository.Video, cr *repository.Comment) *Video {
	return &(Video{userRepo: ur, videoRepo: vr, commentRepo: cr})
}

//服务函数追加在下面
