package service

import (
	"fmt"
	"videohub/internal/repository"
)

type Video_service struct {
	//视频服务，用到user表、video表、comment表操作
	user_repository    *repository.User_repository
	video_repository   *repository.Video_repository
	comment_repository *repository.Comment_repository
}

// 工厂函数，返回单例的服务层操作对象
func NewVideo_service(ur *repository.User_repository, vr *repository.Video_repository, cr *repository.Comment_repository) *Video_service {
	return &(Video_service{user_repository: ur, video_repository: vr, comment_repository: cr})
}

/*
*@author:廖嘉鹏
*@create_at:2024/10/17
 */
// 测试，这是一个样板
func (vs *Video_service) Test() error {
	fmt.Println("Video_service.Test()调用正常")
	vs.user_repository.Test()
	vs.video_repository.Test()
	vs.comment_repository.Test()
	return nil
}

//服务函数追加在下面
