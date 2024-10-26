package service

import (
	"fmt"
	"videohub/internal/repository"
)

type User struct {
	//用户服务，用到user表、collection表、video表操作
	userRepo       *repository.User
	collectionRepo *repository.Collection
	videoRepo      *repository.Video
}

// 工厂函数，返回单例的服务层操作对象
func NewUser(ur *repository.User, cr *repository.Collection, vr *repository.Video) *User {
	return &(User{userRepo: ur, collectionRepo: cr, videoRepo: vr})
}

/*
*@author:廖嘉鹏
*@create_at:2024/10/17
 */
// 测试，这是一个样板
func (us *User) Test() error {
	fmt.Println("User_service.Test()调用正常")
	us.userRepo.Test()
	us.videoRepo.Test()
	us.collectionRepo.Test()
	return nil
}

//服务函数追加在下面
