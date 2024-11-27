package repository

import (
	"videohub/internal/model"

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

// CreateVideo 保存完整视频到数据库
func (vr *Video) CreateVideo(value *model.Video) error {
	return vr.DB.Model(&model.Video{}).Create(value).Error
}

// UpdateVideoStatus 更新视频状态
func (vr *Video) UpdateVideoStatus(id string, newStatus int8) error {
	return vr.DB.Model(&model.Video{}).Where("upload_id = ?", id).Update("video_status", newStatus).Error
}

// 查询视频列表
func (vr *Video) FindVideos(like string, status, page, limit int) ([]model.Video, int64, error) {
	var videos []model.Video
	var count int64

	query := vr.DB.Model(&model.Video{})
	query = query.Where("video_status = ?", status)

	// 标题模糊搜索
	if like != "" {
		query = query.Where("title LIKE ?", "%"+like+"%")
	}

	// 计算偏移量
	offset := (page - 1) * limit
	// 分页查询
	err := query.Count(&count).Offset(offset).Limit(limit).Find(&videos).Error
	if err != nil {
		return nil, 0, err
	}

	return videos, count, nil
}
