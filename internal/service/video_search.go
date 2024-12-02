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
	var response video.GetVideosResponse
	if err := vs.videoRepo.FindVideos(request.Like, *request.Status, request.Page, request.Limit, &response.Videos); err != nil {
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}
	response.Count = int64(len(response.Videos))
	response.Page = request.Page
	response.Limit = request.Limit
	return utils.Ok(http.StatusOK, &response)
}
