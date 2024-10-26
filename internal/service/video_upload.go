package service

import (
	"fmt"
	"videohub/internal/repository"
)

type Video_upload struct {
	//视频上传服务，用到user表、video表操作
	userRepo  *repository.User
	videoRepo *repository.Video
}

// 工厂函数，返回单例的服务层操作对象
func NewVideo_upload(ur *repository.User, vr *repository.Video) *Video_upload {
	return &(Video_upload{userRepo: ur, videoRepo: vr})
}

/*
*@author:廖嘉鹏
*@create_at:2024/10/17
 */
// 测试，这是一个样板
func (vus *Video_upload) Test() error {
	fmt.Println("Video_upload_service.Test()调用正常")
	vus.userRepo.Test()
	vus.videoRepo.Test()
	return nil
}

//服务函数追加在下面
