package service

import (
	"net/http"
	"videohub/internal/repository"
	"videohub/internal/utils"
	"videohub/internal/utils/video"

	"github.com/sirupsen/logrus"
)

// Like 提供点赞业务逻辑
type Like struct {
	videoRepo *repository.Video
	likeRepo  *repository.Like
}

// NewLike 创建点赞服务
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
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}
	err = li.likeRepo.IncrementVideoLikes(request.VideoID)
	if err != nil {
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	logrus.Debug("Like video successfully")
	return utils.Success(http.StatusOK)
}

// UnlikeVideo 取消点赞视频
func (li *Like) UnlikeVideo(request *video.UnLikeVideoRequest) *utils.Response {
	err := li.likeRepo.RemoveVideoLikeRecord(request.UserID, request.VideoID)
	if err != nil {
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}
	err = li.likeRepo.DecrementVideoLikes(request.VideoID)
	if err != nil {
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	logrus.Debug("Unlike video successfully")
	return utils.Success(http.StatusOK)
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
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}
	err = li.likeRepo.IncrementCommentLikes(request.CommentID)
	if err != nil {
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	logrus.Debug("Like comment successfully")
	return utils.Success(http.StatusOK)
}

// UnlikeComment 取消点赞评论
func (li *Like) UnlikeComment(request *video.UnLikeCommentRequest) *utils.Response {
	err := li.likeRepo.RemoveCommentLikeRecord(request.UserID, request.CommentID)
	if err != nil {
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}
	err = li.likeRepo.DecrementCommentLikes(request.CommentID)
	if err != nil {
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	logrus.Debug("Unlike comment successfully")
	return utils.Success(http.StatusOK)
}
