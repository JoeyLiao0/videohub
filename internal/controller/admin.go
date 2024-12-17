package controller

import (
	"net/http"
	"videohub/config"
	"videohub/internal/service"
	"videohub/internal/utils"
	"videohub/internal/utils/admin"
	"videohub/internal/utils/user"
	"videohub/internal/utils/video"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AdminController 包含用户相关的服务
type AdminController struct {
	userAvatarService  *service.UserAvatar
	userListService    *service.UserList
	userService        *service.User
	videoListService   *service.VideoList
	videoUpdateService *service.VideoUpdateStatus
	userVideoService   *service.UserVideo
	statsService       *service.Stats
}

// NewAdminController 创建一个新的 AdminController 实例
func NewAdminController(
	uas *service.UserAvatar,
	uls *service.UserList,
	us *service.User,
	vls *service.VideoList,
	vus *service.VideoUpdateStatus,
	uvs *service.UserVideo,
	d *service.Stats,
) *AdminController {
	return &AdminController{
		userAvatarService:  uas,
		userListService:    uls,
		userService:        us,
		videoListService:   vls,
		videoUpdateService: vus,
		userVideoService:   uvs,
		statsService:       d,
	}
}

// GetUser 获取管理员个人信息
func (ac *AdminController) GetUser(c *gin.Context) {
	id, err := GetUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "请求无效"))
		return
	}
	response := ac.userService.GetUserByID(id)
	c.JSON(http.StatusOK, response)
}

// GetUsers 获取用户信息
func (ac *AdminController) GetUsers(c *gin.Context) {
	var request admin.ListUsersRequest
	if err := c.ShouldBind(&request); err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "请求无效"))
		return
	}

	response := ac.userListService.GetUsers(&request)
	c.JSON(http.StatusOK, response)
}

// CreateUser 创建用户
func (ac *AdminController) CreateUser(c *gin.Context) {
	var request admin.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "请求参数错误"))
		return
	}

	response := ac.userService.CreateUserByAdmin(&request)
	c.JSON(http.StatusOK, response)
}

// UpdateUser 更新用户信息
func (ac *AdminController) UpdateUser(c *gin.Context) {
	var request admin.UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "请求参数错误"))
		return
	}

	response := ac.userService.UpdateUserByAdmin(&request)
	c.JSON(http.StatusOK, response)
}

// GetVideos 获取视频列表
func (ac *AdminController) GetVideos(c *gin.Context) {
	var request admin.GetVideosRequest
	if err := c.ShouldBind(&request); err != nil { // 输入为json
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "请求无效"))
		return
	}

	defaultStatus := -1
	if request.Status == nil {
		request.Status = &defaultStatus
	}

	if request.Page == 0 {
		request.Page = config.AppConfig.Video.DefaultPage
	}

	if request.Limit == 0 {
		request.Limit = config.AppConfig.Video.DefaultLimit
	}

	response := ac.videoListService.GetVideos(&request)
	c.JSON(http.StatusOK, response)
}

// UpdateVideo 更新视频信息
func (ac *AdminController) UpdateVideo(c *gin.Context) {
	// 使用 video 包中的请求结构体
	var request video.UpdateVideoStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "请求无效"))
		return
	}

	response := ac.videoUpdateService.UpdateVideoStatus(&request)
	c.JSON(http.StatusOK, response)
}

// DeleteVideo 删除视频
func (ac *AdminController) DeleteVideo(c *gin.Context) {
	var request user.DeleteVideoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "请求参数错误"))
		return
	}

	// 使用管理员专用的删除方法
	response := ac.userVideoService.DeleteVideoByAdmin(&request)
	c.JSON(http.StatusOK, response)
}

func (ac *AdminController) GetRealTimeData(c *gin.Context) {
	resonse := ac.statsService.GetRealTimeData()
	c.JSON(http.StatusOK, resonse)
}

func (ac *AdminController) GetHistoricalData(c *gin.Context) {
	var request admin.GetHistoricalDataRequest
	if err := c.ShouldBind(&request); err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "请求无效"))
		return
	}
	response := ac.statsService.GetHistoricalData(&request)
	c.JSON(http.StatusOK, response)
}
