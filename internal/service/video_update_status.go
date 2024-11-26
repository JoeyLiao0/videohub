package service

import (
	"errors"
	"log"
	"videohub/internal/repository"
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
func (vus *VideoUpdateStatus) UpdateVideoStatus(id int64, newStatus int8) error {
	// 验证状态合法性
	if newStatus < 0 || newStatus > 3 {
		return errors.New("非法的视频状态")
	}
	if vus.videoRepo == nil {
		log.Printf("videoRepo is nil")
		return errors.New("videoRepo 未初始化")
	}

	log.Printf("VideoRepo is not nil, calling UpdateVideoStatus")
	// 调用DAO更新状态
	return vus.videoRepo.UpdateVideoStatus(id, newStatus)
}
