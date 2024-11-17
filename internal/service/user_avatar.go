package service

import (
	"fmt"
	"mime/multipart"
	"videohub/internal/repository"
)

type User_avatar struct {
	//头像服务，只用到user表操作，所以只用注入user_repostiory
	userRepo *repository.User
}

// 工厂函数，返回单例的服务层操作对象
func NewUser_avatar(ur *repository.User) *User_avatar {
	return &(User_avatar{userRepo: ur})
}

// UploadUserAvatar 上传用户头像
func (uas *User_avatar) UploadUserAvatar(userID uint, file multipart.File) error {
	// 调用 user_repository 进行头像上传相关操作
	// 这里可以是头像文件的存储操作等
	fmt.Printf("Uploading avatar for user ID: %d\n", userID)
	// 实现具体的头像上传逻辑，例如保存文件到本地或云存储
	return nil
}
