package model

import (
	"time"
)

// LikeRecord 点赞记录表
type LikeRecord struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`                              // 主键 ID
	UserID    uint      `gorm:"not null;index:idx_user_target,unique" json:"user_id"`            // 用户 ID
	VideoID   string    `gorm:"size:255;index:idx_user_target,unique" json:"video_id,omitempty"` // 视频 ID，可为空
	CommentID uint      `gorm:"index:idx_user_target,unique" json:"comment_id,omitempty"`        // 评论 ID，可为空
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`                                // 点赞时间
}

// TableName 指定表名为 like_records
func (LikeRecord) TableName() string {
	return "like_records"
}
