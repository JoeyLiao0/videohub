package repository

import (
	"videohub/internal/model"

	"gorm.io/gorm"
)

// Comment 提供评论数据访问接口
type Comment struct {
	DB *gorm.DB
}

func NewComment(db *gorm.DB) *Comment {
	return &Comment{DB: db}
}

// GetCommentsByVideo获取指定视频的评论列表
func (r *Comment) GetCommentsByVideo(videoID int64) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.DB.Where("video_id = ? AND status = 0", videoID).Find(&comments).Error
	return comments, err
}

// CreateComment创建新的评论
func (r *Comment) CreateComment(comment *model.Comment) error {
	return r.DB.Create(comment).Error
}

// DeleteComment删除评论
func (r *Comment) DeleteComment(commentID int64) error {
	// 将所有子评论的status设置为1（标记为已删除）
	if err := r.DB.Model(&model.Comment{}).Where("parent_id = ?", commentID).Update("status", 1).Error; err != nil {
		return err
	}

	return r.DB.Model(&model.Comment{}).Where("id = ?", commentID).Update("status", 1).Error
}
