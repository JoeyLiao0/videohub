package service

import (
	"fmt"
	"net/http"
	"path/filepath"
	"videohub/config"
	"videohub/internal/repository"
	"videohub/internal/utils"
	"videohub/internal/utils/user"

	"github.com/sirupsen/logrus"
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
func (uas *UserAvatar) UploadUserAvatar(id uint, request *user.UploadAvatarRequest) *utils.Response {
	if err := utils.CheckFile(request.Avatar, []string{".png", ".jpg", ".jpeg"}, 8<<20); err != nil {
		logrus.Debug(err.Error())
		return utils.Error(http.StatusBadRequest, "文件格式错误或文件过大")
	}

	fileExt := filepath.Ext(request.Avatar.Filename)
	filePath := filepath.Join(config.AppConfig.Storage.Images, fmt.Sprintf("%d%s", id, fileExt))
	if err := utils.SaveFile(request.Avatar, filePath); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, err.Error())
	}

	if err := uas.userRepo.Update(map[string]interface{}{"id": id}, "avatar", map[string]interface{}{"avatar": config.AppConfig.Storage.Base + "/" + filePath}); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	logrus.Debug("Upload user avatar successfully")
	return utils.Success(http.StatusOK)
}
