package controller

import (
	"net/http"
	"videohub/config"
	"videohub/internal/service"
	"videohub/internal/utils"
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
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "请求无效"))
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
	c.JSON(http.StatusOK, response)
}

// UpdateVideoStatusHandler 处理视频状态更新请求
func (vc *VideoController) UpdateVideoStatus(c *gin.Context) {
	var request video.UpdateVideoStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "请求无效"))
		return
	}

	// 调用服务层更新视频状态
	response := vc.videoUpdateStatus.UpdateVideoStatus(&request)
	c.JSON(http.StatusOK, response)
}

// LikeVideo 点赞视频
func (vc *VideoController) LikeVideo(c *gin.Context) {
	// TODO
}

// GetComments 获取视频评论
func (vc *VideoController) GetComments(c *gin.Context) {
	var request video.GetCommentsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "请求无效"))
		return
	}
	response := vc.comment.GetComments(&request)
	c.JSON(http.StatusOK, response)
}

// AddComment 添加视频评论
func (vc *VideoController) AddComment(c *gin.Context) {
	var request video.AddCommentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "请求无效"))
		return
	}

	response := vc.comment.CreateComment(&request)
	c.JSON(http.StatusOK, response)
}

// LikeComment 点赞评论
func (vc *VideoController) LikeComment(c *gin.Context) {
	// TODO
}

// UnlikeComment 取消点赞评论
func (vc *VideoController) UnlikeComment(c *gin.Context) {
	// TODO
}

// DeleteComment 删除评论
func (vc *VideoController) DeleteComment(c *gin.Context) {
	id, err := GetUserID(c)
	if err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusUnauthorized, "未授权"))
		return
	}

	var request video.DeleteCommentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "请求无效"))
		return
	}

	response := vc.comment.DeleteComment(id, &request)
	c.JSON(http.StatusOK, response)
}

// UploadChunk 处理切片上传请求
func (vc *VideoController) UploadChunk(c *gin.Context) {
	var videoChunk video.UploadChunkRequest
	if err := c.ShouldBind(&videoChunk); err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "请求无效"))
		return
	}

	response := vc.videoUpload.HandleVideoChunk(&videoChunk)
	c.JSON(http.StatusOK, response)
}

// CompleteUpload 处理完整视频合并请求
func (vc *VideoController) CompleteUpload(c *gin.Context) {
	var request video.CompleteUploadRequest
	if err := c.ShouldBind(&request); err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "请求无效"))
		return
	}

	response := vc.videoUpload.HandleVideoComplete(&request)
	c.JSON(http.StatusOK, response)
}
