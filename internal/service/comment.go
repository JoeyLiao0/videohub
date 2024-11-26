package service

import (
	"errors"
	"videohub/internal/model"
	"videohub/internal/repository"
)

// CommentService 提供评论业务逻辑
type CommentService struct {
	Repo *repository.Comment
}

func NewCommentService(repo *repository.Comment) *CommentService {
	return &CommentService{Repo: repo}
}

// GetComments获取视频的所有评论
func (s *CommentService) GetComments(videoID int64) ([]model.Comment, error) {
	return s.Repo.GetCommentsByVideo(videoID)
}

// CreateComment创建评论s
func (s *CommentService) CreateComment(userID int64, videoID int64, content string, parentID int64) error {
	comment := &model.Comment{
		UserID:         userID,
		VideoID:        videoID,
		CommentContent: content,
		ParentID:       parentID,
		Status:         0,
	}
	return s.Repo.CreateComment(comment)
}

// DeleteComment删除评论
func (s *CommentService) DeleteComment(commentID int64) error {
	// TODO: 权限？
	err := s.Repo.DeleteComment(commentID)
	if err != nil {
		return errors.New("删除评论失败")
	}
	return nil
}
