package controller

import (
	"net/http"
	"strconv"
	"videohub/internal/model"
	"videohub/internal/service"

	"github.com/gin-gonic/gin"
)

type VideosController struct {
	videoUpload *service.VideoUpload
}

// NewVideoController creates a new instance of VideosController with the provided service.
func NewVideoController(vus *service.VideoUpload) *VideosController {
	return &VideosController{videoUpload: vus}
}

// GetVideos 获取视频列表
func (vc *VideosController) GetVideos(c *gin.Context) {
	// TODO
}

// LikeVideo 点赞视频
func (vc *VideosController) LikeVideo(c *gin.Context) {
	// TODO
}

// GetComments 获取视频评论
func (vc *VideosController) GetComments(c *gin.Context) {
	// TODO
}

// AddComment 添加视频评论
func (vc *VideosController) AddComment(c *gin.Context) {
	// TODO
}

// LikeComment 点赞评论
func (vc *VideosController) LikeComment(c *gin.Context) {
	// TODO
}

// DeleteComment 删除评论
func (vc *VideosController) DeleteComment(c *gin.Context) {
	// TODO
}

// UploadChunk 处理切片上传请求
func (vc *VideosController) UploadChunk(c *gin.Context) {
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
func (vc *VideosController) CompleteUpload(c *gin.Context) {
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
