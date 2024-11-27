package controller

import (
	"net/http"
	"strconv"
	"videohub/config"
	"videohub/internal/service"
	"videohub/internal/utils/video"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type VideoController struct {
	videoUpload       *service.VideoUpload
	videoUpdateStatus *service.VideoUpdateStatus
	videoSearch       *service.VideoSearch
	comment           *service.Comment
}

func NewVideoController(videoUpload *service.VideoUpload, videoUpdateStatus *service.VideoUpdateStatus, videoSearch *service.VideoSearch, comment *service.Comment) *VideoController {
	return &VideoController{
		videoUpload:       videoUpload,
		videoUpdateStatus: videoUpdateStatus,
		videoSearch:       videoSearch,
		comment:           comment,
	}
}

// GetVideos 获取视频列表
func (vc *VideoController) GetVideos(c *gin.Context) {
	// 获取 Query 参数
	// 或者使用 c.DefaultQuery()
	var request video.GetVideosRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的查询参数"})
		return
	}

	if request.Status == nil {
		request.Status = &config.AppConfig.Video.DefaultStatus
	}

	if request.Page == 0 {
		request.Page = config.AppConfig.Video.DefaultPage
	}

	if request.Limit == 0 {
		request.Limit = config.AppConfig.Video.DefaultLimit
	}

	// 调用服务层获取视频列表
	response := vc.videoSearch.GetVideos(&request)
	c.JSON(response.StatusCode, response)
}

// UpdateVideoStatusHandler 处理视频状态更新请求
func (vc *VideoController) UpdateVideoStatus(c *gin.Context) {
	// 获取视频 ID
	vid := c.Param("vid")

	var request video.UpdateVideoStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析请求体"})
		return
	}

	// 调用服务层更新视频状态
	response := vc.videoUpdateStatus.UpdateVideoStatus(vid, &request)
	c.JSON(response.StatusCode, response)
}

// LikeVideo 点赞视频
func (vc *VideoController) LikeVideo(c *gin.Context) {
	// TODO
}

// GetComments 获取视频评论
func (vc *VideoController) GetComments(c *gin.Context) {
	vid := c.Param("vid")
	response := vc.comment.GetComments(vid)
	c.JSON(response.StatusCode, response)
}

// AddComment 添加视频评论
func (vc *VideoController) AddComment(c *gin.Context) {
	vid := c.Param("vid")
	var request video.AddCommentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析请求体"})
		return
	}

	response := vc.comment.CreateComment(vid, &request)
	c.JSON(response.StatusCode, response)
}

// LikeComment 点赞评论
func (vc *VideoController) LikeComment(c *gin.Context) {
	// TODO
}

// DeleteComment 删除评论
func (vc *VideoController) DeleteComment(c *gin.Context) {
	id, err := GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的操作"})
		return
	}

	vid := c.Param("vid")
	cid, err := strconv.Atoi(c.Param("cid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的评论 ID"})
		return
	}

	response := vc.comment.DeleteComment(id, vid, uint(cid))
	c.JSON(response.StatusCode, response)
}

// UploadChunk 处理切片上传请求
func (vc *VideoController) UploadChunk(c *gin.Context) {
	var videoChunk video.UploadChunkRequest
	if err := c.ShouldBind(&videoChunk); err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求的Form格式无效或缺少必需字段"})
		return
	}

	response := vc.videoUpload.HandleVideoChunk(&videoChunk)
	c.JSON(response.StatusCode, response)
}

// CompleteUpload 处理完整视频合并请求
func (vc *VideoController) CompleteUpload(c *gin.Context) {
	var request video.CompleteUploadRequest
	if err := c.ShouldBind(&request); err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求的Form格式无效或缺少必需字段"})
		return
	}

	response := vc.videoUpload.HandleVideoComplete(&request)
	c.JSON(response.StatusCode, response)
}
