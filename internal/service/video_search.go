package service

import (
	"videohub/internal/model"
	"videohub/internal/repository"
)

// VideoService 提供视频业务逻辑
type VideoSearch struct {
	videoRepo *repository.Video
}

// VideoService 实例
func NewVideoSearch(vr *repository.Video) *VideoSearch {
	return &VideoSearch{videoRepo: vr}
}

// 获取视频列表
func (vs *VideoSearch) GetVideos(status *int, like *string, page, limit int) ([]model.Video, int64, error) {
	// 计算分页偏移量
	offset := (page - 1) * limit

	// 调用数据层查询
	return vs.videoRepo.FindVideos(status, like, offset, limit)
}
