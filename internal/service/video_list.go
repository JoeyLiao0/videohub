package service

import (
	"net/http"
	"videohub/global"
	"videohub/internal/repository"
	"videohub/internal/utils"
	"videohub/internal/utils/admin"
	"videohub/internal/utils/video"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type VideoList struct {
	//视频列表服务，用到user表、video表操作
	userRepo  *repository.User
	videoRepo *repository.Video
}

// 工厂函数，返回单例的服务层操作对象
func NewVideoList(ur *repository.User, vr *repository.Video) *VideoList {
	return &(VideoList{userRepo: ur, videoRepo: vr})
}

// 服务函数追加在下面
func (vl *VideoList) GetVideos(request *admin.GetVideosRequest) *utils.Response {
	var videos []video.VideoInfo
	if err := vl.videoRepo.FindVideos(
		request.Like,    // 标题模糊搜索
		*request.Status, // 视频状态过滤
		request.Page,    // 当前页码
		request.Limit,   // 每页数量
		[]string{ // 需要查询的字段
			"videos.upload_id",
			"videos.created_at",
			"videos.title",
			"videos.description",
			"videos.cover_path",
			"videos.video_path",
			"videos.video_status",
			"users.username as uploader_name",
			"videos.likes",
			"videos.favorites",
			"videos.comments",
		},
		&videos,
	); err != nil {
		logrus.Error("Failed to find videos: ", err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	for i := range videos {
		videoKey := "video:" + videos[i].UploadID + ":views"
		views, err := global.Rdb.Get(global.Ctx, videoKey).Int()
		if err == redis.Nil {
			videos[i].Views = 0
		} else if err != nil {
			logrus.Error("Failed to get video views: ", err.Error())
			return utils.Error(http.StatusInternalServerError, "服务器内部错误")
		} else {
			videos[i].Views = views
		}
	}

	return utils.Ok(http.StatusOK, &video.GetVideosResponse{
		Videos: videos,
	})
}
