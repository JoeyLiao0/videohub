package video

import "videohub/internal/model"

type VideoInfo struct {
	UploadID       string `json:"id"`
	CreatedAt      int64  `json:"published_at"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	CoverPath      string `json:"cover_path"`
	VideoPath      string `json:"video_path"`
	VideoStatus    int8   `json:"status"`
	UploaderID     int    `json:"-"`
	UploaderName   string `json:"name"`
	UploaderAvatar string `json:"avatar"`
	Likes          int    `json:"like_count"`
	Favorites      int    `json:"collection_count"`
	Comments       int    `json:"comment_count"`
	Views          int    `json:"view_count"`
	IsLiked        bool   `json:"is_liked"`
	IsCollected    bool   `json:"is_collected"`
}

type CommentInfo struct {
	ID             uint   `json:"comment_id"`
	CreatedAt      int64  `json:"created_at"`
	UserID         int    `json:"user_id"`
	Username       string `json:"name"`
	Avatar         string `json:"avatar"`
	CommentContent string `json:"comment"`
	VideoID        string `json:"video_id"`
	ParentID       int    `json:"father_comment_id"`
	Likes          int    `json:"likes_count"`
	Status         int    `json:"-"`
}

type GetVideosResponse struct {
	Videos []VideoInfo `json:"videos"`
}

type CommentsInside struct {
	Comments CommentInfo `json:"comments"`
	IsLiked  bool        `json:"is_liked"`
	ReplyTo  string      `json:"reply_to"`
}
type CommentsOutside struct {
	Comments CommentInfo      `json:"comments"`
	IsLiked  bool             `json:"is_liked"`
	Reply    []CommentsInside `json:"reply"`
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
