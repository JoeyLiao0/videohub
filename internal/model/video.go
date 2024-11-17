package model

import (
	"gorm.io/gorm"
)

type VideoChunk struct {
	UploadID  string `json:"upload_id" binding:"required"`  // 上传任务的唯一ID
	ChunkID   int    `json:"chunk_id" binding:"required"`   // 分片编号
	ChunkSize int    `json:"chunk_size" binding:"required"` // 切片大小
	ChunkHash string `json:"chunk_hash" binding:"required"` // 切片哈希
}

type Video struct {
	gorm.Model
	UploadID    int64  `gorm:"primary_key" json:"upload_id"` // 上传任务的唯一ID
	Title       string `gorm:"text" json:"title"`            // 视频标题
	Description string `gorm:"text" json:"description"`      // 视频描述
	CoverPath   string `gorm:"size:1024" json:"cover_path"`  // 视频封面路径
	VideoPath   string `gorm:"size:1024" json:"video_path"`  // 完整视频的文件路径
	VideoStatus int8   `json:"video_status"`                 // 视频状态（0-正常 1-待审核 2-审核未通过 3-封禁）
	UploaderID  int64  `json:"uploarder_id"`                 // 发布者id
	Likes       int    `json:"likes"`                        // 点赞数
	Favorites   int    `json:"favorites"`                    // 收藏数
	Comments    int    `json:"comments"`                     // 评论数
}
