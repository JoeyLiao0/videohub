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

// UserCollection 用户收藏服务层操作对象
type UserCollection struct {
	videoRepo      *repository.Video
	likeRepo       *repository.Like
	collectionRepo *repository.Collection
}

// NewUserCollection 实例化用户收藏服务层操作对象
func NewUserCollection(vr *repository.Video, lr *repository.Like, cr *repository.Collection) *UserCollection {
	return &UserCollection{videoRepo: vr, likeRepo: lr, collectionRepo: cr}
}

// GetUserCollections 获取用户收藏列表
func (uc *UserCollection) GetUserCollections(id uint) *utils.Response {
	var response user.VideoListResponse
	conditions := map[string]interface{}{"collections.user_id": id}
	joins := []string{"left join videos on collections.video_id = videos.upload_id", "left join users on collections.user_id = users.id"}
	fields := []string{
		"upload_id",
		"created_at",
		"title",
		"description",
		"cover_path",
		"video_path",
		"video_status",
		"likes",
		"favorites",
		"comments",
	}
	for i, field := range fields {
		fields[i] = "videos." + field
	}
	fields = append(fields, "users.username as uploader_name")
	fields = append(fields, "users.avatar as uploader_avatar")
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

		isLiked, err := uc.likeRepo.CheckVideoLike(id, response.Videos[i].UploadID)
		if err != nil {
			logrus.Error(err.Error())
			return utils.Error(http.StatusInternalServerError, "获取视频点赞状态失败")
		}
		response.Videos[i].IsLiked = isLiked
		isCollected, err := uc.collectionRepo.CheckVideoCollect(id, response.Videos[i].UploadID)
		if err != nil {
			logrus.Error(err.Error())
			return utils.Error(http.StatusInternalServerError, "获取视频收藏状态失败")
		}
		response.Videos[i].IsCollected = isCollected
	}

	logrus.Debug("Get user videos successfully")
	return utils.Ok(http.StatusOK, &response)
}

// AddUserCollection 添加用户收藏
func (uc *UserCollection) AddUserCollection(userID uint, request *user.AddCollectionsRequest) *utils.Response {
	// 检查是否已经收藏
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
	err := uc.collectionRepo.IncrementVideoCollects(collection.VideoID)
	if err != nil {
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	logrus.Debug("Add user collection successfully")
	return utils.Success(http.StatusOK)
}

// DeleteUserCollection 删除用户收藏
func (uc *UserCollection) DeleteUserCollection(userID uint, request *user.DeleteCollectionsRequest) *utils.Response {
	if err := uc.collectionRepo.Delete(map[string]interface{}{"user_id": userID, "video_id": request.VideoID}); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	err := uc.collectionRepo.DecrementVideoCollects(request.VideoID)
	if err != nil {
		return utils.Error(http.StatusInternalServerError, "更新视频收藏数失败")
	}
	
	logrus.Debug("Delete user collection successfully")
	return utils.Success(http.StatusOK)
}
