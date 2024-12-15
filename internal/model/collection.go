package model

type Collection struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint   `gorm:"not null" json:"user_id"`  // 收藏者ID
	VideoID   string `gorm:"not null" json:"video_id"` // 被收藏的视频ID
	CreatedAt int64  `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt int64  `gorm:"autoUpdateTime:milli" json:"updated_at"`
}