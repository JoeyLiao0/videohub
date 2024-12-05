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
}

// VideoService 实例
func NewVideoSearch(vr *repository.Video) *VideoSearch {
	return &VideoSearch{videoRepo: vr}
}

// 获取视频列表
func (vs *VideoSearch) GetVideos(request *video.GetVideosRequest) *utils.Response {
	var response video.GetVideosResponse
	var fileds = []string{"upload_id", "created_at", "title", "description", "cover_path", "video_path", "video_status", "uploader_name", "likes", "favorites", "comments"}
	if err := vs.videoRepo.FindVideos(request.Like, *request.Status, request.Page, request.Limit, fileds, &response.Videos); err != nil {
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
	}
	// response.Count = int64(len(response.Videos))
	// response.Page = request.Page
	// response.Limit = request.Limit
	return utils.Ok(http.StatusOK, &response)
}
