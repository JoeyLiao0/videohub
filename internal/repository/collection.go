package repository

import (
	"videohub/internal/model"

	"gorm.io/gorm"
)

// Collection 视频收藏表的数据库操作封装
type Collection struct {
	DB *gorm.DB
}

// NewCollection 实例化视频收藏表的数据库操作
func NewCollection(db *gorm.DB) *Collection {
	return &Collection{DB: db}
}

// Create 创建收藏记录
func (r *Collection) Create(value *model.Collection) error {
	return r.DB.Model(&model.Collection{}).Create(value).Error
}

// Count 统计符合条件的收藏记录数量
func (r *Collection) Count(conditions interface{}) (int64, error) {
	var count int64
	err := r.DB.Model(&model.Collection{}).Where(conditions).Count(&count).Error
	return count, err
}

// Delete 删除收藏记录
func (r *Collection) Delete(conditions interface{}) error {
	return r.DB.Where(conditions).Delete(&model.Collection{}).Error
}

// GetUserCollections 获取用户收藏记录
func (r *Collection) GetUserCollections(conditions interface{}, limit int, joins []string,
	fields, result interface{}) error {
	query := r.DB.Model(&model.Collection{}).Where(conditions).Limit(limit).Select(fields)
	for _, join := range joins {
		query = query.Joins(join)
	}
	return query.Find(result).Error
}

// IncrementVideoCollects 增加视频收藏数
func (co *Collection) IncrementVideoCollects(videoID string) error {
	return co.DB.Model(&model.Video{}).
		Where("upload_id = ?", videoID).
		Update("favorites", gorm.Expr("favorites + ?", 1)).Error
}

// DecrementVideoCollects 减少视频收藏数
func (co *Collection) DecrementVideoCollects(videoID string) error {
	return co.DB.Model(&model.Video{}).
		Where("upload_id = ?", videoID).
		Update("favorites", gorm.Expr("favorites - ?", 1)).Error
}

// CheckVideoCollect 检查是否已经收藏
func (co *Collection) CheckVideoCollect(userID uint, videoID string) (bool, error) {
	var count int64
	err := co.DB.Model(&model.Collection{}).
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Count(&count).Error
	return count > 0, err
}
