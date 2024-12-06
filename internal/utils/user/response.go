package user

import (
	"videohub/internal/utils/video"
)

type UserInfo struct {
	Username  string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Avatar    string `json:"avatar" binding:"required"`
	Status    int8   `json:"status" binding:"required"`
	CreatedAt int64  `json:"time" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token" binding:"required"`
}

type GetUserResponse struct {
	User UserInfo `json:"user" binding:"required"`
}

type VideoListResponse struct {
	Videos []video.VideoInfo `json:"videos"`
}

type GetCommentsResponse struct {
	Comments []video.CommentInfo `json:"comments"`
}
