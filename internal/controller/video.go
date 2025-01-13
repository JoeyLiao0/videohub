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

// VideoController 视频控制器
type VideoController struct {
	videoUpload       *service.VideoUpload       // 视频上传服务
	videoUpdateStatus *service.VideoUpdateStatus // 视频状态更新服务
	videoSearch       *service.VideoSearch       // 视频搜索服务
	like              *service.Like              // 点赞服务
	comment           *service.Comment           // 评论服务
}

// NewVideoController 创建一个新的 VideoController 实例
func NewVideoController(videoUpload *service.VideoUpload, videoUpdateStatus *service.VideoUpdateStatus,
	videoSearch *service.VideoSearch, like *service.Like, comment *service.Comment) *VideoController {
	return &VideoController{
		videoUpload:       videoUpload,
		videoUpdateStatus: videoUpdateStatus,
		videoSearch:       videoSearch,
		like:              like,
		comment:           comment,
	}
}

// GetVideos 获取视频列表
func (vc *VideoController) GetVideos(c *gin.Context) {
	var request video.GetVideosRequest
	if err := c.ShouldBind(&request); err != nil {
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
	// JWT（可有可无）
	token := c.GetHeader("Authorization")
	if token == "" {
		request.UserID = 0
	}
	payload, err := utils.ParseJWT(token, config.AppConfig.JWT.AccessTokenSecret)

	if err != nil {
		request.UserID = 0
	} else {
		request.UserID = payload.ID
	}

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

	response := vc.videoUpdateStatus.UpdateVideoStatus(&request)
	c.JSON(http.StatusOK, response)
}

// LikeVideo 点赞视频
func (vc *VideoController) LikeVideo(c *gin.Context) {
	var request video.LikeVideoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "无效的请求参数"))
		return
	}

	if request.VideoID == "" {
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "视频ID不能为空"))
		return
	}

	userID, _ := GetUserID(c)
	request.UserID = userID

	response := vc.like.LikeVideo(&request)
	c.JSON(http.StatusOK, response)
}

// UnlikeVideo 取消点赞视频
func (vc *VideoController) UnlikeVideo(c *gin.Context) {
	var request video.UnLikeVideoRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "无效的请求参数"))
		return
	}

	if request.VideoID == "" {
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "视频ID不能为空"))
		return
	}

	userID, _ := GetUserID(c)
	request.UserID = userID

	response := vc.like.UnlikeVideo(&request)
	c.JSON(http.StatusOK, response)
}

// GetComments 获取视频评论
func (vc *VideoController) GetComments(c *gin.Context) {
	var request video.GetCommentsRequest
	if err := c.ShouldBind(&request); err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "无效的请求参数"))
		return
	}

	// JWT（可有可无）
	token := c.GetHeader("Authorization")
	if token == "" {
		request.UserID = 0
	}
	payload, err := utils.ParseJWT(token, config.AppConfig.JWT.AccessTokenSecret)
	if err != nil {
		request.UserID = 0
	} else {
		request.UserID = payload.ID
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

	userID, _ := GetUserID(c)
	request.UserID = userID

	response := vc.comment.CreateComment(&request)
	c.JSON(http.StatusOK, response)
}

// LikeComment 点赞评论
func (vc *VideoController) LikeComment(c *gin.Context) {
	var request video.LikeCommentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "无效的请求参数"))
		return
	}

	userID, _ := GetUserID(c)
	request.UserID = userID

	response := vc.like.LikeComment(&request)
	c.JSON(http.StatusOK, response)
}

// UnlikeComment 取消点赞评论
func (vc *VideoController) UnlikeComment(c *gin.Context) {
	var request video.UnLikeCommentRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "无效的请求参数"))
		return
	}
	userID, _ := GetUserID(c)
	request.UserID = userID

	response := vc.like.UnlikeComment(&request)
	c.JSON(http.StatusOK, response)
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
	if err := c.ShouldBind(&request); err != nil {
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
	id, err := GetUserID(c)
	if err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusUnauthorized, "未授权"))
		return
	}
	var request video.CompleteUploadRequest
	if err := c.ShouldBind(&request); err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "请求无效"))
		return
	}

	response := vc.videoUpload.HandleVideoComplete(id, &request)
	c.JSON(http.StatusOK, response)
}
