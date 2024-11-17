package controller

import (
	"net/http"
	"strconv"
	"videohub/internal/model"
	"videohub/internal/service"

	"github.com/gin-gonic/gin"
)

type Videos struct {
	videoUpload *service.VideoUpload
}

func NewVideos(vus *service.VideoUpload) *Videos {
	return &Videos{videoUpload: vus}
}

// UploadChunk 处理切片上传请求
func (controller *Videos) UploadChunk(c *gin.Context) {
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

	if err := controller.videoUpload.HandleVideoChunk(videoChunk, fileHeader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "chunk processing error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Chunk uploaded successfully"})
}

// CompleteUpload 处理完整视频合并请求
func (controller *Videos) CompleteUpload(c *gin.Context) {
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

	if err := controller.videoUpload.HandleVideoComplete(video, chunkEndID, coverFile, videoHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "video merge error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Video uploaded successfully"})
}
