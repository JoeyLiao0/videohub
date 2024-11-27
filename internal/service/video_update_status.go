package service

import (
	"net/http"
	"videohub/internal/repository"
	"videohub/internal/utils"
	"videohub/internal/utils/video"

	"github.com/sirupsen/logrus"
)

// VideoService 提供视频业务逻辑
type VideoUpdateStatus struct {
	videoRepo *repository.Video
}

// VideoService实例
func NewVideoUpdateStatus(vr *repository.Video) *VideoUpdateStatus {
	return &VideoUpdateStatus{videoRepo: vr}
}

// UpdateVideoStatus更新视频状态
func (vus *VideoUpdateStatus) UpdateVideoStatus(id string, request *video.UpdateVideoStatusRequest) *utils.Response {
	// 验证状态合法性
	if request.NewStatus < 0 || request.NewStatus > 3 {
		return utils.Error(http.StatusBadRequest, "无效的视频状态")
	}

	if err := vus.videoRepo.UpdateVideoStatus(id, request.NewStatus); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "更新视频状态失败")
	}

	return utils.Success(http.StatusOK)
}
