package video

import "videohub/internal/model"

type GetVideosResponse struct {
	Videos []model.Video `json:"videos"`
	Page   int           `json:"page"`
	Limit  int           `json:"limit"`
	Count  int64         `json:"count"`
}

type GetCommentsResponse struct {
	Comments []model.Comment `json:"comments"`
}