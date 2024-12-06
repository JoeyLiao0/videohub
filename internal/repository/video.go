package repository

import (
	"videohub/global"
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
func (vr *Video) FindVideos(like string, status, page, limit int, fields, result interface{}) error {
	query := vr.DB.Model(&model.Video{})
	query = query.Where("videos.video_status = ?", status)

	// 标题模糊搜索
	if like != "" {
		query = query.Where("videos.title LIKE ?", "%"+like+"%")
	}

	// 计算偏移量
	offset := (page - 1) * limit
	// 分页查询
	err := query.Offset(offset).Limit(limit).Select(fields).Joins("left join users on videos.uploader_id = users.id").Find(result).Error
	if err != nil {
		return err
	}
	return nil
}
