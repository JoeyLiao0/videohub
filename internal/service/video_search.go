package service

import (
	"net/http"
	"videohub/global"
	"videohub/internal/repository"
	"videohub/internal/utils"
	"videohub/internal/utils/video"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// VideoService 提供视频业务逻辑
type VideoSearch struct {
	videoRepo *repository.Video
	likeRepo  *repository.Like
}

// VideoService 实例
func NewVideoSearch(vr *repository.Video, lr *repository.Like) *VideoSearch {
	return &VideoSearch{videoRepo: vr, likeRepo: lr}
}

// 获取视频列表
func (vs *VideoSearch) GetVideos(request *video.GetVideosRequest) *utils.Response {
	var response video.GetVideosResponse
	fields := []string{"upload_id", "created_at", "title", "description", "cover_path", "video_path", "video_status", "likes", "favorites", "comments"}
	for i, field := range fields {
		fields[i] = "videos." + field
	}
	fields = append(fields, "users.username as uploader_name")
	if err := vs.videoRepo.FindVideos(request.Like, *request.Status, request.Page, request.Limit, fields, &response.Videos); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	for i := range response.Videos {
		views, err := global.Rdb.Get(global.Ctx, "video:"+response.Videos[i].UploadID+":views").Int()
		if err == redis.Nil {
			logrus.Debug("redis: nil")
			views = 0
		} else if err != nil {
			logrus.Debug(err.Error())
			return utils.Error(http.StatusInternalServerError, "服务器内部错误")
		}
		response.Videos[i].Views = views
		// 查询用户是否点赞
		isLiked, err := vs.likeRepo.CheckVideoLike(request.UserID, response.Videos[i].UploadID)
		if err != nil {
			logrus.Debug(err.Error())
			return utils.Error(http.StatusInternalServerError, "服务器内部错误")
		}
		response.Videos[i].IsLiked = isLiked
	}

	return utils.Ok(http.StatusOK, &response)
}
