package repository

import (
	"videohub/global"
	"videohub/internal/model"
	"videohub/internal/utils/video"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Video struct {
	DB *gorm.DB
}

func NewVideo(db *gorm.DB) *Video {
	return &Video{DB: db}
}

func (vr *Video) Search(conditions interface{}, limit int, result interface{}) error {
	return vr.DB.Model(&model.Video{}).Where(conditions).Limit(limit).Find(result).Error
}

func (vr *Video) Select(conditions interface{}, limit int, fields, result interface{}) error {
	return vr.DB.Model(&model.Video{}).Where(conditions).Limit(limit).Select(fields).Find(result).Error
}

func (vr *Video) Join(conditions interface{}, limit int, joins string, fields, result interface{}) error {
	return vr.DB.Model(&model.Video{}).Where(conditions).Limit(limit).Select(fields).Joins(joins).Find(result).Error
}

func (vr *Video) Delete(conditions interface{}) error {
	return vr.DB.Where(conditions).Delete(&model.Video{}).Error
}

// CreateVideo 保存完整视频到数据库
func (vr *Video) CreateVideo(value *model.Video) error {
	global.Rdb.Set(global.Ctx, "video:"+value.UploadID+":views", 0, 0)
	return vr.DB.Model(&model.Video{}).Create(value).Error
}

// UpdateVideoStatus 更新视频状态
func (vr *Video) UpdateVideoStatus(id string, newStatus int8) error {
	return vr.DB.Model(&model.Video{}).Where("upload_id = ?", id).Update("video_status", newStatus).Error
}

// 查询视频列表
func (vr *Video) GetVideos(like string, status, page, limit int) ([]video.VideoInfo, error) {
	var videoInfos []video.VideoInfo

	// 查询字段
	fields := []string{
		"videos.upload_id",
		"videos.created_at",
		"videos.title",
		"videos.description",
		"videos.cover_path",
		"videos.video_path",
		"videos.video_status",
		"videos.likes",
		"videos.favorites",
		"videos.comments",
		"videos.uploader_id",
	}

	// 构建查询
	offset := (page - 1) * limit
	query := vr.DB.Debug().Model(&model.Video{}).
		Select(fields).
		Where("videos.video_status = ?", status).
		Offset(offset).
		Limit(limit) // 偏移分页

	// 标题模糊搜索
	if like != "" {
		query = query.Where("videos.title LIKE ?", "%"+like+"%")
	}

	// 执行查询
	if err := query.Scan(&videoInfos).Error; err != nil {
		return nil, err
	}

	// 填充观看数（从服务层移到这里）
	for i := range videoInfos {
		views, err := global.Rdb.Get(global.Ctx, "video:"+videoInfos[i].UploadID+":views").Int()
		if err == redis.Nil {
			views = 0
		} else if err != nil {
			logrus.Error(err.Error())
			return nil, err
		}
		videoInfos[i].Views = views

		// 手动查询填充username和avatar
		var user model.User
		if err := vr.DB.Debug().Model(&model.User{}).
			Where("id = ?", videoInfos[i].UploaderID). // 使用 uploader_id
			First(&user).Error; err != nil {
			logrus.Warnf("User not found for uploader_id: %d", videoInfos[i].UploaderID)
		} else {
			videoInfos[i].UploaderName = user.Username
			videoInfos[i].UploadAvatar = user.Avatar
		}
	}

	return videoInfos, nil
}
