package service

import (
	"net/http"
	"videohub/global"
	"videohub/internal/model"
	"videohub/internal/repository"
	"videohub/internal/utils"
	"videohub/internal/utils/user"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type UserCollection struct {
	collectionRepo *repository.Collection
}

func NewUserCollection(cr *repository.Collection) *UserCollection {
	return &UserCollection{collectionRepo: cr}
}

func (uc *UserCollection) GetUserCollections(id uint) *utils.Response {
	var response user.VideoListResponse
	conditions := map[string]interface{}{"collections.user_id": id}
	joins := []string{"left join videos on collections.video_id = videos.upload_id", "left join users on collections.user_id = users.id"}
	fields := []string{"upload_id", "created_at", "title", "description", "cover_path", "video_path", "video_status", "likes", "favorites", "comments"}
	for i, field := range fields {
		fields[i] = "videos." + field
	}
	fields = append(fields, "users.username as uploader_name")
	if err := uc.collectionRepo.GetUserCollections(conditions, -1, joins, fields, &response.Videos); err != nil {
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
	logrus.Debug("Get user collections successfully")
	return utils.Ok(http.StatusOK, &response)
}

func (uc *UserCollection) AddUserCollection(userID uint, request *user.AddCollectionsRequest) *utils.Response {
	if count, err := uc.collectionRepo.Count(map[string]interface{}{"user_id": userID, "video_id": request.VideoID}); err != nil || count > 0 {
		if err != nil {
			logrus.Error(err.Error())
			return utils.Error(http.StatusInternalServerError, "服务器内部错误")
		} else {
			logrus.Debug("Video has been collected")
			return utils.Error(http.StatusBadRequest, "视频已收藏")
		}
	}

	collection := &model.Collection{
		UserID:  userID,
		VideoID: request.VideoID,
	}
	if err := uc.collectionRepo.Create(collection); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	logrus.Debug("Add user collection successfully")
	return utils.Ok(http.StatusOK, nil)
}

func (uc *UserCollection) DeleteUserCollection(userID uint, request *user.DeleteCollectionsRequest) *utils.Response {
	if err := uc.collectionRepo.Delete(map[string]interface{}{"user_id": userID, "video_id": request.VideoID}); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	logrus.Debug("Delete user collection successfully")
	return utils.Ok(http.StatusOK, nil)
}
