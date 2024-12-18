package service

import (
	"net/http"
	"videohub/internal/repository"
	"videohub/internal/utils"
	"videohub/internal/utils/video"

	"github.com/sirupsen/logrus"
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
	// 检查是否已经点赞
	isLiked, err := li.likeRepo.CheckVideoLike(request.UserID, request.VideoID)
	if isLiked {
		if err != nil {
			logrus.Error(err.Error())
			return utils.Error(http.StatusInternalServerError, "服务器内部错误")
		} else {
			logrus.Debug("Video has been liked")
			return utils.Error(http.StatusBadRequest, "视频已点赞")
		}
	}
	err = li.likeRepo.AddVideoLikeRecord(request.UserID, request.VideoID)
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
	// 检查是否已经点赞
	isLiked, err := li.likeRepo.CheckCommentLike(request.UserID, request.CommentID)
	if isLiked {
		if err != nil {
			logrus.Error(err.Error())
			return utils.Error(http.StatusInternalServerError, "服务器内部错误")
		} else {
			logrus.Debug("Comment has been liked")
			return utils.Error(http.StatusBadRequest, "评论已点赞")
		}
	}
	err = li.likeRepo.AddCommentLikeRecord(request.UserID, request.CommentID)
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
