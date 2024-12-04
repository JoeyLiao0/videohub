package model

type Comment struct {
	ID             uint   `json:"id" gorm:"primaryKey;autoIncrement"`     // 评论ID
	CreatedAt      int64  `gorm:"autoCreateTime:milli" json:"created_at"` // 评论发布时间
	Username       string `json:"name" gorm:"not null"`                   // 评论发布者ID
	CommentContent string `json:"comment"`                                // 评论内容
	VideoID        string `json:"video_id" gorm:"not null"`               // 视频唯一标识
	ParentID       int    `json:"father_comment_id"`                      // 父评论ID
	Likes          int    `json:"likes"`                                  // 评论点赞数
	Status         int    `json:"status"`                                 // 评论状态（0-正常，1-删除）
}
