package service

import (
	"net/http"
	"videohub/internal/repository"
	"videohub/internal/utils"
	"videohub/internal/utils/video"
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
func (vs *VideoSearch) GetVideos(request *video.GetVideosRequest) *utils.Response {
	// 调用数据层查询
	videos, count, err := vs.videoRepo.FindVideos(request.Like, *request.Status, request.Page, request.Limit)
	if err != nil {
		return utils.Error(http.StatusInternalServerError, "视频列表查询失败")
	}
	return utils.Ok(http.StatusOK, &video.GetVideosResponse{
		Videos: videos,
		Page:   request.Page,
		Limit:  request.Limit,
		Count:  count,
	})
}
