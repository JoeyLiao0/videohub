package service

import (
	"fmt"
	"videohub/internal/repository"
)

type Video_list_service struct {
	//视频列表服务，用到user表、video表操作
	user_repository  *repository.User_repository
	video_repository *repository.Video_repository
}

// 工厂函数，返回单例的服务层操作对象
func NewVideo_list_service(ur *repository.User_repository, vr *repository.Video_repository) *Video_list_service {
	return &(Video_list_service{user_repository: ur, video_repository: vr})
}

/*
*@author:廖嘉鹏
*@create_at:2024/10/17
 */
// 测试，这是一个样板
func (vls *Video_list_service) Test() error {
	fmt.Println("Video_list_service.Test()调用正常")
	vls.user_repository.Test()
	vls.video_repository.Test()
	return nil
}

//服务函数追加在下面
