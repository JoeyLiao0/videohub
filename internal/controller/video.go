package controller

import (
	"net/http"
	"strconv"
	"videohub/internal/model"
	"videohub/internal/service"

	"github.com/gin-gonic/gin"
)

type VideoController struct {
	videoUpload       *service.VideoUpload
	videoUpdateStatus *service.VideoUpdateStatus
	videoSearch       *service.VideoSearch
	comment           *service.CommentService
}

func NewVideoController(videoUpload *service.VideoUpload, videoUpdateStatus *service.VideoUpdateStatus, videoSearch *service.VideoSearch, comment *service.CommentService) *VideoController {
	return &VideoController{
		videoUpload:       videoUpload,
		videoUpdateStatus: videoUpdateStatus,
		videoSearch:       videoSearch,
		comment:           comment,
	}
}

// GetVideos 获取视频列表
func (vc *VideoController) GetVideos(c *gin.Context) {
	// 获取Query参数
	var (
		status *int
		like   *string
		page   = 1
		limit  = 10
	)

	if statusStr := c.Query("status"); statusStr != "" {
		statusVal, err := strconv.Atoi(statusStr)
		if err == nil {
			status = &statusVal
		}
	}

	if likeStr := c.Query("like"); likeStr != "" {
		like = &likeStr
	}

	if pageStr := c.Query("page"); pageStr != "" {
		pageVal, err := strconv.Atoi(pageStr)
		if err == nil {
			page = pageVal
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		limitVal, err := strconv.Atoi(limitStr)
		if err == nil {
			limit = limitVal
		}
	}

	// 调用服务层获取视频列表
	videos, total, err := vc.videoSearch.GetVideos(status, like, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "服务器错误", "error": err.Error()})
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"videos": videos,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		},
	})
}

// UpdateVideoStatusHandler 处理视频状态更新请求
func (vc *VideoController) UpdateVideoStatus(c *gin.Context) {
	// 获取视频ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的视频 ID"})
		return
	}

	// 解析JSON
	var requestBody struct {
		NewStatus int8 `json:"new_status"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析请求体"})
		return
	}

	// 调用服务层更新视频状态
	err = vc.videoUpdateStatus.UpdateVideoStatus(id, requestBody.NewStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"msg": "视频状态更新成功"})
}

// LikeVideo 点赞视频
func (vc *VideoController) LikeVideo(c *gin.Context) {
	// TODO
}

// GetComments 获取视频评论
func (vc *VideoController) GetComments(c *gin.Context) {
	videoID, err := strconv.ParseInt(c.Param("vid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	comments, err := vc.comment.GetComments(videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// AddComment 添加视频评论
func (vc *VideoController) AddComment(c *gin.Context) {
	videoID, err := strconv.ParseInt(c.Param("vid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	var req struct {
		UserID         int64  `json:"user_id"`
		CommentContent string `json:"comment"`
		ParentID       int64  `json:"father_comment_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = vc.comment.CreateComment(req.UserID, videoID, req.CommentContent, req.ParentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment created successfully"})
}

// LikeComment 点赞评论
func (vc *VideoController) LikeComment(c *gin.Context) {
	// TODO
}

// DeleteComment 删除评论
func (vc *VideoController) DeleteComment(c *gin.Context) {
	_, err := strconv.ParseInt(c.Param("vid"), 10, 64) // vid?
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	commentID, err := strconv.ParseInt(c.Param("cid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	err = vc.comment.DeleteComment(commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// UploadChunk 处理切片上传请求
func (vc *VideoController) UploadChunk(c *gin.Context) {
	var videoChunk model.VideoChunk

	if err := c.ShouldBindJSON(&videoChunk); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 从form-data中获取文件
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file upload error"})
		return
	}

	if err := vc.videoUpload.HandleVideoChunk(videoChunk, fileHeader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "chunk processing error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Chunk uploaded successfully"})
}

// CompleteUpload 处理完整视频合并请求
func (vc *VideoController) CompleteUpload(c *gin.Context) {
	// 从表单数据获取参数
	uploadID, _ := strconv.ParseInt(c.PostForm("upload_id"), 10, 64)
	title := c.PostForm("title")
	description := c.PostForm("description")
	uploaderID, _ := strconv.ParseInt(c.PostForm("uploader_id"), 10, 64)
	videoHash := c.PostForm("video_hash")

	// Video结构体实例，方便后续向服务层函数传入参数
	video := model.Video{
		UploadID:    uploadID,
		Title:       title,
		Description: description,
		UploaderID:  uploaderID,
	}

	// 从form-data中获取文件
	coverFileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "coverFile upload error"})
		return
	}

	// 打开封面文件
	coverFile, err := coverFileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to open cover file"})
		return
	}
	defer coverFile.Close()

	chunkEndID, _ := strconv.Atoi(c.PostForm("chunk_end_id"))

	if err := vc.videoUpload.HandleVideoComplete(video, chunkEndID, coverFile, videoHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "video merge error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Video uploaded successfully"})
}
