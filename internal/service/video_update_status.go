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
func (vus *VideoUpdateStatus) UpdateVideoStatus(request *video.UpdateVideoStatusRequest) *utils.Response {
	status := *request.NewStatus
	// 验证状态合法性
	if status < 0 || status > 3 {
		return utils.Error(http.StatusBadRequest, "无效的视频状态")
	}

	if err := vus.videoRepo.UpdateVideoStatus(request.VideoID, status); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "更新视频状态失败")
	}

	logrus.Debug("Video status updated successfully")
	return utils.Success(http.StatusOK)
}
