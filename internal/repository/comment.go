package repository

import (
	"log"
	"videohub/internal/model"
	"videohub/internal/utils/video"

	"gorm.io/gorm"
)

// Comment 提供评论数据访问接口
type Comment struct {
	DB *gorm.DB
}

func NewComment(db *gorm.DB) *Comment {
	return &Comment{DB: db}
}

func (r *Comment) Search(conditions interface{}, limit int, result interface{}) error {
	return r.DB.Model(&model.Comment{}).Where(conditions).Limit(limit).Find(result).Error
}

func (r *Comment) Select(conditions interface{}, limit int, fields, result interface{}) error {
	return r.DB.Model(&model.Comment{}).Where(conditions).Limit(limit).Select(fields).Find(result).Error
}

func (r *Comment) Join(conditions interface{}, limit int, joins string, fields, result interface{}) error {
	return r.DB.Model(&model.Comment{}).Where(conditions).Limit(limit).Select(fields).Joins(joins).Find(result).Error
}

// GetCommentsByVideo获取指定视频的评论列表
func (r *Comment) GetCommentsByVideo(videoID string, uid uint) ([]video.CommentsOutside, error) {
	var parentComments []model.Comment
	// 查询父评论（ParentID == -1 且 Status == 0）
	err := r.DB.Where("video_id = ? AND status = 0 AND parent_id = -1", videoID).Find(&parentComments).Error
	if err != nil {
		return nil, err
	}

	var commentsOutside []video.CommentsOutside
	for _, parentComment := range parentComments {
		var isLiked bool
		// 填充外部评论的字段
		if uid == 0 {
			isLiked = false
		} else {
			isLiked, err = r.CheckIsLiked(parentComment.ID, uid) // uid应该为jwt对应的uid
		}
		if err != nil {
			return nil, err
		}

		commentsOutside = append(commentsOutside, video.CommentsOutside{
			Comments: parentComment,
			IsLiked:  isLiked,
			Reply:    nil, // 子评论稍后填充
		})
	}

	return commentsOutside, nil
}

// FillPerCommentsReply 填充子评论并处理IsLiked和ReplyTo字段
func (r *Comment) FillPerCommentsReply(commentsOutside []video.CommentsOutside, uid uint) ([]video.CommentsOutside, error) {
	for i, comment := range commentsOutside {
		// 使用递归函数填充所有子评论
		replies, err := r.getRepliesRecursive(comment.Comments.ID, comment.Comments.UserID, uid)
		if err != nil {
			return nil, err
		}

		// 更新父评论的回复字段
		commentsOutside[i].Reply = replies
	}

	return commentsOutside, nil
}

// getRepliesRecursive 递归获取子评论，并填充ReplyTo字段
func (r *Comment) getRepliesRecursive(parentID uint, parentUserID uint, uid uint) ([]video.CommentsInside, error) {
	var childComments []model.Comment
	// 查询直接子评论
	err := r.DB.Where("parent_id = ? AND status = 0", parentID).Find(&childComments).Error
	if err != nil {
		return nil, err
	}

	var replies []video.CommentsInside
	for _, childComment := range childComments {
		// 填充子评论的点赞状态
		var isLiked bool
		if uid == 0 {
			isLiked = false
		} else {
			isLiked, err = r.CheckIsLiked(childComment.ID, uid)
		}
		if err != nil {
			return nil, err
		}

		// 获取父评论用户的名称
		replyToName := r.GetUserName(parentUserID)

		// 递归获取子评论的子评论（孙子评论）
		grandReplies, err := r.getRepliesRecursive(childComment.ID, childComment.UserID, uid)
		if err != nil {
			return nil, err
		}

		// 当前评论
		currentReply := video.CommentsInside{
			Comments: childComment,
			IsLiked:  isLiked,
			ReplyTo:  replyToName,
		}

		// 添加当前评论到回复列表
		replies = append(replies, currentReply)

		// 添加孙子评论到回复列表
		replies = append(replies, grandReplies...)
	}

	return replies, nil
}

// GetUserName获取reply的用户名
func (r *Comment) GetUserName(userID uint) string {
	var user model.User
	err := r.DB.Select("username").Where("id = ?", userID).First(&user).Error
	if err != nil {
		// 如果查询失败，返回空字符串或默认名称
		log.Printf("Error retrieving username for userID %d: %v", userID, err)
		return "查询失败"
	}
	return user.Username
}

// CheckIsLiked 判断是否点赞
func (r *Comment) CheckIsLiked(commentID uint, userID uint) (bool, error) {
	var likeRecord model.LikeRecord
	err := r.DB.Where("comment_id = ? AND user_id = ?", commentID, userID).First(&likeRecord).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CreateComment创建新的评论
func (r *Comment) CreateComment(comment *model.Comment) error {
	return r.DB.Create(comment).Error
}

// DeleteComment删除评论
func (r *Comment) DeleteComment(cid uint) error {
	// 将所有子评论的status设置为1（标记为已删除）
	if err := r.DB.Model(&model.Comment{}).Where("parent_id = ?", cid).Update("status", 1).Error; err != nil {
		return err
	}

	return r.DB.Model(&model.Comment{}).Where("id = ?", cid).Update("status", 1).Error
}
