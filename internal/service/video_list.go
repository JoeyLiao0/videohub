package service

import (
	"fmt"
	"videohub/internal/repository"
)

type Video_list struct {
	//视频列表服务，用到user表、video表操作
	userRepo  *repository.User
	videoRepo *repository.Video
}

// 工厂函数，返回单例的服务层操作对象
func NewVideo_list(ur *repository.User, vr *repository.Video) *Video_list {
	return &(Video_list{userRepo: ur, videoRepo: vr})
}

/*
*@author:廖嘉鹏
*@create_at:2024/10/17
 */
// 测试，这是一个样板
func (vls *Video_list) Test() error {
	fmt.Println("Video_list_service.Test()调用正常")
	vls.userRepo.Test()
	vls.videoRepo.Test()
	return nil
}

//服务函数追加在下面
