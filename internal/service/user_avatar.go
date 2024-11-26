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
	if err := utils.CheckFile(request.Avatar, []string{".png", ".jpg"}, 8<<20); err != nil {
		logrus.Debug(err.Error())
		return utils.Error(http.StatusBadRequest, err.Error())
	}

	fileExt := filepath.Ext(request.Avatar.Filename)
	filePath := fmt.Sprintf("%s/%d%s", config.AppConfig.Storage.Images, id, fileExt)
	if err := utils.SaveFile(request.Avatar, filePath); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, err.Error())
	}

	if err := uas.userRepo.Update(map[string]interface{}{"id": id}, "avatar", map[string]interface{}{"avatar": filePath}); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "更新用户头像信息失败")
	}

	logrus.Debug("Upload user avatar successfully")
	return utils.Success(http.StatusOK)
}
