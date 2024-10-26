package service

import (
	"fmt"
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

/*
*@author:廖嘉鹏
*@create_at:2024/10/17
 */
// 测试，这是一个样板
func (vs *Video) Test() error {
	fmt.Println("Video_service.Test()调用正常")
	vs.userRepo.Test()
	vs.videoRepo.Test()
	vs.commentRepo.Test()
	return nil
}

//服务函数追加在下面
