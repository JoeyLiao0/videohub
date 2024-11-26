package model

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	ID             int64     `json:"id" gorm:"primaryKey;autoIncrement"` // 评论ID
	UserID         int64     `json:"user_id" gorm:"not null"`            // 评论发布者ID
	CommentContent string    `json:"comment"`                            // 评论内容
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`   // 评论发布时间
	VideoID        int64     `json:"video_id" gorm:"not null"`           // 视频唯一标识
	ParentID       int64     `json:"father_comment_id"`                  // 父评论ID
	Likes          int       `json:"likes"`                              // 评论点赞数
	Status         int       `json:"status"`                             // 评论状态（0-正常，1-删除）
}