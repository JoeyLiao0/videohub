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

// UserAvatar 用户头像服务层操作对象
type UserAvatar struct {
	userRepo *repository.User
}

// NewUserAvatar 实例化用户头像服务层操作对象
func NewUserAvatar(ur *repository.User) *UserAvatar {
	return &(UserAvatar{userRepo: ur})
}

// UploadUserAvatar 上传用户头像
func (uas *UserAvatar) UploadUserAvatar(id uint, request *user.UploadAvatarRequest) *utils.Response {
	if err := utils.CheckFile(request.Avatar, []string{".png", ".jpg", ".jpeg"}, 8<<20); err != nil {
		logrus.Debug(err.Error())
		return utils.Error(http.StatusBadRequest, "文件格式错误或文件过大")
	}

	var result struct {
		Avatar string
	}

	if err := uas.userRepo.Search(map[string]interface{}{"id": id}, 1, &result); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	if filepath.Base(result.Avatar) != "tourist.png" {
		if err := utils.RemoveFile(filepath.Join(config.AppConfig.Storage.Images, filepath.Base(result.Avatar))); err != nil {
			logrus.Error(err.Error())
			return utils.Error(http.StatusInternalServerError, "服务器内部错误")
		}
	}

	fileExt := filepath.Ext(request.Avatar.Filename)
	filePath := filepath.Join(config.AppConfig.Storage.Images, fmt.Sprintf("%d%s", id, fileExt))
	if err := utils.SaveFile(request.Avatar, filePath); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, err.Error())
	}

	values := map[string]interface{}{
		"avatar": utils.GetURLPath(config.AppConfig.Static.Avatar, fmt.Sprintf("%d%s", id, fileExt)),
	}
	if err := uas.userRepo.Update(map[string]interface{}{"id": id}, "avatar", values); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	logrus.Debug("Upload user avatar successfully")
	return utils.Success(http.StatusOK)
}
