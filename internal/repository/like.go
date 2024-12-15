package repository

import (
	"videohub/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Like struct {
	DB *gorm.DB
}

func NewLike(db *gorm.DB) *Like {
	return &Like{DB: db}
}

// CheckVideoLike 检查用户是否已经点赞过视频
func (li *Like) CheckVideoLike(userID uint, videoID string) (bool, error) {
	var count int64
	err := li.DB.Table("like_records").
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Count(&count).Error
	logrus.Debugf("CheckVideoLike - userID: %d, videoID: %s, count: %d", userID, videoID, count)
	return count > 0, err
}

// CheckCommentLike 检查用户是否已经点赞过评论
func (li *Like) CheckCommentLike(userID uint, commentID uint) (bool, error) {
	var count int64
	err := li.DB.Table("like_records").
		Where("user_id = ? AND comment_id = ?", userID, commentID).
		Count(&count).Error
	return count > 0, err
}

// AddVideoLikeRecord 添加视频点赞记录
func (li *Like) AddVideoLikeRecord(userID uint, videoID string) error {
	record := map[string]interface{}{
		"user_id":    userID,
		"video_id":   videoID,
		"created_at": gorm.Expr("NOW()"),
	}
	return li.DB.Table("like_records").Create(&record).Error
}

// AddCommentLikeRecord 添加评论点赞记录
func (li *Like) AddCommentLikeRecord(userID uint, commentID uint) error {
	record := map[string]interface{}{
		"user_id":    userID,
		"comment_id": commentID,
		"created_at": gorm.Expr("NOW()"),
	}
	return li.DB.Table("like_records").Create(&record).Error
}

// IncrementVideoLikes 增加视频点赞数
func (li *Like) IncrementVideoLikes(videoID string) error {
	return li.DB.Model(&model.Video{}).
		Where("upload_id = ?", videoID).
		Update("likes", gorm.Expr("likes + ?", 1)).Error
}

// IncrementCommentLikes 增加评论点赞数
func (li *Like) IncrementCommentLikes(commentID uint) error {
	return li.DB.Model(&model.Comment{}).
		Where("id = ?", commentID).
		Update("likes", gorm.Expr("likes + ?", 1)).Error
}

// RemoveVideoLikeRecord 删除视频点赞记录
func (li *Like) RemoveVideoLikeRecord(userID uint, videoID string) error {
	return li.DB.Table("like_records").
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Delete(nil).Error
}

// DecrementVideoLikes 减少视频点赞数
func (li *Like) DecrementVideoLikes(videoID string) error {
	return li.DB.Model(&model.Video{}).
		Where("upload_id = ?", videoID).
		Update("likes", gorm.Expr("likes - ?", 1)).Error
}

// RemoveCommentLikeRecord 删除评论点赞记录
func (li *Like) RemoveCommentLikeRecord(userID uint, commentID uint) error {
	return li.DB.Table("like_records").
		Where("user_id = ? AND comment_id = ?", userID, commentID).
		Delete(nil).Error
}

// DecrementCommentLikes 减少评论点赞数
func (li *Like) DecrementCommentLikes(commentID uint) error {
	return li.DB.Model(&model.Comment{}).
		Where("id = ?", commentID).
		Update("likes", gorm.Expr("likes - ?", 1)).Error
}
