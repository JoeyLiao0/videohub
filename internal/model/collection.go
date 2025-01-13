package model

// Collection 视频收藏表
type Collection struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`     // 收藏 ID
	UserID    uint   `gorm:"not null" json:"user_id"`                // 收藏者 ID
	VideoID   string `gorm:"not null" json:"video_id"`               // 被收藏的视频 ID
	CreatedAt int64  `gorm:"autoCreateTime:milli" json:"created_at"` // 收藏时间
	UpdatedAt int64  `gorm:"autoUpdateTime:milli" json:"updated_at"` // 更新时间
}
