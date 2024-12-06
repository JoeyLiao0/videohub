package service

import (
	"net/http"
	"videohub/global"
	"videohub/internal/repository"
	"videohub/internal/utils"
	"videohub/internal/utils/user"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type UserVideo struct {
	videoRepo *repository.Video // 视频仓库的指针，用于操作视频数据
}

// NewUserVideo 创建一个新的 UserVideo 实例
func NewUserVideo(vr *repository.Video) *UserVideo {
	return &UserVideo{videoRepo: vr}
}

func (uv *UserVideo) GetUserVideos(id uint) *utils.Response {
	var response user.VideoListResponse
	conditions := map[string]interface{}{"videos.uploader_id": id}
	joins := "left join users on videos.uploader_id = users.id"
	fields := []string{"upload_id", "created_at", "title", "description", "cover_path", "video_path", "video_status", "likes", "favorites", "comments"}
	for i, field := range fields {
		fields[i] = "videos." + field
	}
	fields = append(fields, "users.username as uploader_name")
	if err := uv.videoRepo.Join(conditions, -1, joins, fields, &response.Videos); err != nil {
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
	}
	logrus.Debug("Get user videos successfully")
	return utils.Ok(http.StatusOK, &response)
}

func (uv *UserVideo) DeleteUserVideo(id uint, request *user.DeleteVideoRequest) *utils.Response {
	conditions := map[string]interface{}{
		"uploader_id": id,
		"upload_id": request.VideoID,
	}
	if err := uv.videoRepo.Delete(conditions); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}
	logrus.Debug("Delete video successfully")
	return utils.Success(http.StatusOK)
}