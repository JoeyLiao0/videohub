package video

type VideoInfo struct {
	UploadID    string `json:"upload_id"` // 上传任务的唯一ID
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	Title       string `json:"title"`        // 视频标题
	Description string `json:"description"`  // 视频描述
	CoverPath   string `json:"cover_path"`   // 视频封面路径
	VideoPath   string `json:"video_path"`   // 完整视频的文件路径
	VideoStatus int8   `json:"video_status"` // 视频状态（0-正常 1-待审核 2-审核未通过 3-封禁）
	// TODO: 应该修改为视频上传者的用户名
	UploaderID  uint   `json:"uploarder_id"` // 发布者id
	Likes       int    `json:"likes"`        // 点赞数
	Favorites   int    `json:"favorites"`    // 收藏数
	Comments    int    `json:"comments"`     // 评论数
}

type CommentInfo struct {
	ID        uint  `json:"id"`         // 评论ID
	CreatedAt int64 `json:"created_at"` // 评论发布时间
	// TODO: 应该修改为评论发布者的用户名
	UserID         uint   `json:"user_id"`           // 评论发布者ID
	CommentContent string `json:"comment"`           // 评论内容
	VideoID        string `json:"video_id"`          // 视频唯一标识
	ParentID       int    `json:"father_comment_id"` // 父评论ID
	Likes          int    `json:"likes"`             // 评论点赞数
	Status         int    `json:"status"`            // 评论状态（0-正常，1-删除）
}

type GetVideosResponse struct {
	Videos []VideoInfo `json:"videos"`
	Page   int         `json:"page"`
	Limit  int         `json:"limit"`
	Count  int64       `json:"count"`
}

type GetCommentsResponse struct {
	Comments []CommentInfo `json:"comments"`
}
