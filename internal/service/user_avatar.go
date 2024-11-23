package service

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"videohub/config"
	"videohub/internal/repository"
	"videohub/internal/utils"
)

type UserAvatar struct {
	//头像服务，只用到user表操作，所以只用注入user_repostiory
	userRepo *repository.User
}

// 工厂函数，返回单例的服务层操作对象
func NewUserAvatar(ur *repository.User) *UserAvatar {
	return &(UserAvatar{userRepo: ur})
}

// UploadUserAvatar 上传用户头像
func (uas *UserAvatar) UploadUserAvatar(id uint64, file *multipart.FileHeader) *utils.Response {
	fileExt := filepath.Ext(file.Filename)
	filePath := fmt.Sprintf("%s/%d%s", config.AppConfig.Storage.Images, id, fileExt)
	if err := utils.SaveFile(file, filePath); err != nil {
		return utils.Error(http.StatusInternalServerError, err.Error())
	}
	if err := uas.userRepo.Update(map[string]interface{}{"id": id}, map[string]interface{}{"avatar": filePath}); err != nil {
		return utils.Error(http.StatusInternalServerError, "更新用户头像信息失败")
	}

	return utils.Success(http.StatusOK)
}
