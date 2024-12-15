package video

import "videohub/internal/model"

type VideoInfo struct {
	UploadID     string `json:"id"`               // 视频 ID
	CreatedAt    int64  `json:"published_at"`     // 视频发布时间
	Title        string `json:"title"`            // 视频标题
	Description  string `json:"description"`      // 视频描述
	CoverPath    string `json:"cover_path"`       // 视频封面路径
	VideoPath    string `json:"video_path"`       // 完整视频的文件路径
	VideoStatus  int8   `json:"status"`           // 视频状态（0-正常 1-待审核 2-审核未通过 3-封禁）
	UploaderName string `json:"name"`             // 发布者用户名
	Likes        int    `json:"like_count"`       // 点赞数
	Favorites    int    `json:"collection_count"` // 收藏数
	Comments     int    `json:"comment_count"`    // 评论数
	Views        int    `json:"view_count"`       // 观看数
	IsLiked      bool   // 是否点赞
}

type CommentInfo struct {
	ID             uint   `json:"id"`                // 评论ID
	CreatedAt      int64  `json:"created_at"`        // 评论发布时间
	Username       string `json:"name"`              // 评论发布者用户名
	CommentContent string `json:"comment"`           // 评论内容
	VideoID        string `json:"video_id"`          // 视频唯一标识
	ParentID       int    `json:"father_comment_id"` // 父评论ID
	Likes          int    `json:"likes"`             // 评论点赞数
	Status         int    `json:"status"`            // 评论状态（0-正常，1-删除）
}

type GetVideosResponse struct {
	Videos []VideoInfo `json:"videos"`
	// Page   int         `json:"page"`
	// Limit  int         `json:"limit"`
	// Count  int64       `json:"count"`
}

type CommentsInside struct {
	Comments model.Comment `json:"comments"`
	IsLiked  bool          `json:"is_liked"`
	ReplyTo  string        `json:"reply_to"`
}
type CommentsOutside struct {
	Comments model.Comment    `json:"comments"`
	IsLiked  bool             `json:"is_liked"`
	Reply    []CommentsInside `json:"reply"`
}

type GetCommentsResponse struct {
	CommentsOutside []CommentsOutside `json:"comments"`
}
