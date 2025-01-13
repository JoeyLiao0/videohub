package service

import (
	"math"
	"net/http"
	"videohub/internal/repository"
	"videohub/internal/utils"
	"videohub/internal/utils/admin"
	"videohub/internal/utils/video"

	"github.com/sirupsen/logrus"
)

// VideoService 提供视频业务逻辑
type VideoSearch struct {
	videoRepo      *repository.Video
	likeRepo       *repository.Like
	collectionRepo *repository.Collection
}

// NewVideoSearch 实例化视频搜索服务
func NewVideoSearch(vr *repository.Video, lr *repository.Like, cr *repository.Collection) *VideoSearch {
	return &VideoSearch{videoRepo: vr, likeRepo: lr, collectionRepo: cr}
}

// 获取视频列表
func (vs *VideoSearch) GetVideos(request *video.GetVideosRequest) *utils.Response {
	// 计算总记录数
	total, err := vs.videoRepo.Count(*request.Status, request.Like)
	if err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}
	totalPages := int(math.Ceil(float64(total) / float64(request.Limit)))

	videos, err := vs.videoRepo.GetVideos(request.Like, *request.Status, request.Page, request.Limit)
	if err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	// 检查是否点赞、收藏
	for i := range videos {
		if request.UserID != 0 {
			isLiked, err := vs.likeRepo.CheckVideoLike(request.UserID, videos[i].UploadID)
			if err != nil {
				logrus.Error(err.Error())
				return utils.Error(http.StatusInternalServerError, "服务器内部错误")
			}
			videos[i].IsLiked = isLiked
			isCollected, err := vs.collectionRepo.CheckVideoCollect(request.UserID, videos[i].UploadID)
			if err != nil {
				logrus.Error(err.Error())
				return utils.Error(http.StatusInternalServerError, "服务器内部错误")
			}
			videos[i].IsCollected = isCollected
		}
	}

	response := admin.VideoInfo{
		Videos: videos,
		Pages: admin.PageInfo{
			Page:       request.Page,
			Limit:      request.Limit,
			TotalPages: totalPages,
		},
	}

	logrus.Debug("Get videos successfully")
	return utils.Ok(http.StatusOK, &response)
}
