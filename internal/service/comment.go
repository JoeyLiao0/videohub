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
func (cs *Comment) GetComments(vid string) *utils.Response {
	comments, err := cs.commentRepo.GetCommentsByVideo(vid)
	if err != nil {
		return utils.Error(http.StatusInternalServerError, "获取评论失败")
	}
	return utils.Ok(http.StatusOK, &video.GetCommentsResponse{Comments: comments})
}

// CreateComment创建评论s
func (cs *Comment) CreateComment(vid string, request *video.AddCommentRequest) *utils.Response {
	comment := &model.Comment{
		UserID:         request.UserID,
		VideoID:        vid,
		CommentContent: request.CommentContent,
		ParentID:       request.FatherCommentID,
		Status:         0,
	}
	if err := cs.commentRepo.CreateComment(comment); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "评论失败")
	}

	return utils.Success(http.StatusOK)
}

// DeleteComment删除评论
func (cs *Comment) DeleteComment(id uint, vid string, cid uint) *utils.Response {
	// up 和 评论发布者可以删除该评论
	var result struct {
		UserID uint
	}
	if err := cs.commentRepo.Search(map[string]interface{}{"id": cid}, 1, &result); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "删除评论失败")
	}

	var result2 struct {
		UploaderID uint
	}
	if err := cs.videoRepo.Search(map[string]interface{}{"upload_id": vid}, 1, &result2); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "删除评论失败")
	}

	if result.UserID != id && result2.UploaderID != id {
		return utils.Error(http.StatusUnauthorized, "未授权的操作")
	}

	if err := cs.commentRepo.DeleteComment(cid); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "删除评论失败")
	}

	return utils.Success(http.StatusOK)
}
