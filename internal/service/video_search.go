package service

import (
	"net/http"
	"videohub/internal/repository"
	"videohub/internal/utils"
	"videohub/internal/utils/video"

	"github.com/sirupsen/logrus"
)

// VideoService 提供视频业务逻辑
type VideoSearch struct {
	videoRepo      *repository.Video
	likeRepo       *repository.Like
	collectionRepo *repository.Collection
}

// VideoService 实例
func NewVideoSearch(vr *repository.Video, lr *repository.Like, cr *repository.Collection) *VideoSearch {
	return &VideoSearch{videoRepo: vr, likeRepo: lr, collectionRepo: cr}
}

// 获取视频列表
func (vs *VideoSearch) GetVideos(request *video.GetVideosRequest) *utils.Response {
	var response video.GetVideosResponse

	// 从数据层获取视频信息
	videos, err := vs.videoRepo.GetVideos(request.Like, *request.Status, request.Page, request.Limit)
	if err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "获取视频列表失败")
	}

	// 检查是否点赞、收藏
	for i := range videos {
		if request.UserID != 0 {
			isLiked, err := vs.likeRepo.CheckVideoLike(request.UserID, videos[i].UploadID)
			if err != nil {
				logrus.Error(err.Error())
				return utils.Error(http.StatusInternalServerError, "获取视频点赞状态失败")
			}
			videos[i].IsLiked = isLiked
			isCollected, err := vs.collectionRepo.CheckVideoCollect(request.UserID, videos[i].UploadID)
			if err != nil {
				logrus.Error(err.Error())
				return utils.Error(http.StatusInternalServerError, "获取视频收藏状态失败")
			}
			videos[i].IsCollected = isCollected
		}
	}

	response.Videos = videos
	return utils.Ok(http.StatusOK, &response)
}
