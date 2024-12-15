package service

import (
	"videohub/internal/repository"
	"videohub/internal/utils"
	"videohub/internal/utils/video"
)

type Like struct {
	videoRepo *repository.Video
	likeRepo  *repository.Like
}

func NewLike(vr *repository.Video, lr *repository.Like) *Like {
	return &Like{
		videoRepo: vr,
		likeRepo:  lr,
	}
}

// LikeVideo 点赞视频
func (li *Like) LikeVideo(request *video.LikeVideoRequest) *utils.Response {
	err := li.likeRepo.AddVideoLikeRecord(request.UserID, request.VideoID)
	if err != nil {
		return utils.Error(500, "添加点赞记录失败")
	}
	err = li.likeRepo.IncrementVideoLikes(request.VideoID)
	if err != nil {
		return utils.Error(500, "更新视频点赞数失败")
	}
	return utils.Ok(200, "视频点赞成功")
}

// UnlikeVideo 取消点赞视频
func (li *Like) UnlikeVideo(request *video.LikeVideoRequest) *utils.Response {
	err := li.likeRepo.RemoveVideoLikeRecord(request.UserID, request.VideoID)
	if err != nil {
		return utils.Error(500, "删除点赞记录失败")
	}
	err = li.likeRepo.DecrementVideoLikes(request.VideoID)
	if err != nil {
		return utils.Error(500, "更新视频点赞数失败")
	}
	return utils.Ok(200, "视频取消点赞成功")
}

// LikeComment 点赞评论
func (li *Like) LikeComment(request *video.LikeCommentRequest) *utils.Response {
	err := li.likeRepo.AddCommentLikeRecord(request.UserID, request.CommentID)
	if err != nil {
		return utils.Error(500, "添加点赞记录失败")
	}
	err = li.likeRepo.IncrementCommentLikes(request.CommentID)
	if err != nil {
		return utils.Error(500, "更新评论点赞数失败")
	}
	return utils.Ok(200, "评论点赞成功")
}

// UnlikeComment 取消点赞评论
func (li *Like) UnlikeComment(request *video.LikeCommentRequest) *utils.Response {
	err := li.likeRepo.RemoveCommentLikeRecord(request.UserID, request.CommentID)
	if err != nil {
		return utils.Error(500, "删除点赞记录失败")
	}
	err = li.likeRepo.DecrementCommentLikes(request.CommentID)
	if err != nil {
		return utils.Error(500, "更新评论点赞数失败")
	}
	return utils.Ok(200, "评论取消点赞成功")
}
