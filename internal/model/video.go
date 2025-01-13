package model

// Video 视频表
type Video struct {
	UploadID    string `gorm:"primary_key" json:"upload_id"`           // 视频唯一标识
	CreatedAt   int64  `gorm:"autoCreateTime:milli" json:"created_at"` // 视频创建时间
	UpdatedAt   int64  `gorm:"autoUpdateTime:milli" json:"updated_at"` // 视频更新时间
	Title       string `gorm:"text" json:"title"`                      // 视频标题
	Description string `gorm:"text" json:"description"`                // 视频描述
	CoverPath   string `gorm:"size:255" json:"cover_path"`             // 视频封面路径
	VideoPath   string `gorm:"size:255" json:"video_path"`             // 完整视频的文件路径
	VideoStatus int8   `gorm:"defalut:1" json:"video_status"`          // 视频状态: 0-正常 1-待审核 2-审核未通过 3-封禁
	UploaderID  uint   `json:"uploader_id"`                            // 发布者 ID
	Likes       int    `json:"likes"`                                  // 点赞数
	Favorites   int    `json:"favorites"`                              // 收藏数
	Comments    int    `json:"comments"`                               // 评论数
}
