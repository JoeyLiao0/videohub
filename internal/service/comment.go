package service

import (
	"net/http"
	"videohub/internal/model"
	"videohub/internal/repository"
	"videohub/internal/utils"
	"videohub/internal/utils/video"

	"github.com/sirupsen/logrus"
)

// Comment 提供评论业务逻辑
type Comment struct {
	commentRepo *repository.Comment
	videoRepo   *repository.Video
}

func NewComment(cr *repository.Comment, vr *repository.Video) *Comment {
	return &Comment{commentRepo: cr, videoRepo: vr}
}

// GetComments获取视频的所有评论
func (cs *Comment) GetComments(request *video.GetCommentsRequest) *utils.Response {
	vid := request.VideoID
	uid := request.UserID
	// uid 为 0 的时候点赞全部为 false

	// comments中先放入父评论为-1的评论，每个评论的reply数组先为空
	comments, err := cs.commentRepo.GetCommentsByVideo(vid, uid)
	if err != nil {
		return utils.Error(http.StatusInternalServerError, "获取评论失败")
	}
	// 填充每个评论的reply数组
	comments, err = cs.commentRepo.FillPerCommentsReply(comments, uid)
	if err != nil {
		return utils.Error(http.StatusInternalServerError, "获取子评论失败")
	}
	return utils.Ok(http.StatusOK, &video.GetCommentsResponse{CommentsOutside: comments})
}

// CreateComment创建评论
func (cs *Comment) CreateComment(request *video.AddCommentRequest) *utils.Response {
	comment := &model.Comment{
		UserID:         request.UserID,
		VideoID:        request.VideoID,
		CommentContent: request.CommentContent,
		ParentID:       request.FatherCommentID,
		Status:         0,
	}
	if err := cs.commentRepo.CreateComment(comment); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	return utils.Success(http.StatusOK)
}

// DeleteComment删除评论
func (cs *Comment) DeleteComment(id uint, request *video.DeleteCommentRequest) *utils.Response {
	// up 和 评论发布者可以删除该评论
	var result struct {
		UserID uint
	}
	if err := cs.commentRepo.Search(map[string]interface{}{"id": request.CommentID}, 1, &result); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	var result2 struct {
		UploaderID uint
	}
	if err := cs.videoRepo.Search(map[string]interface{}{"upload_id": request.VideoID}, 1, &result2); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	if result.UserID != id && result2.UploaderID != id {
		return utils.Error(http.StatusUnauthorized, "未授权")
	}

	if err := cs.commentRepo.DeleteComment(uint(request.CommentID)); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	return utils.Success(http.StatusOK)
}
