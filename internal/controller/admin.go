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

// AdminController 管理员控制器
type AdminController struct {
	userAvatarService  *service.UserAvatar        // 用户头像服务
	userListService    *service.UserList          // 用户列表服务
	userService        *service.User              // 用户服务
	videoSearchService *service.VideoSearch       // 视频搜索服务
	videoUpdateService *service.VideoUpdateStatus // 视频更新服务
	userVideoService   *service.UserVideo         // 用户视频服务
	statsService       *service.Stats             // 统计服务
}

// NewAdminController 创建一个新的 AdminController 实例
func NewAdminController(
	uas *service.UserAvatar,
	uls *service.UserList,
	us *service.User,
	vls *service.VideoSearch,
	vus *service.VideoUpdateStatus,
	uvs *service.UserVideo,
	d *service.Stats,
) *AdminController {
	return &AdminController{
		userAvatarService:  uas,
		userListService:    uls,
		userService:        us,
		videoSearchService: vls,
		videoUpdateService: vus,
		userVideoService:   uvs,
		statsService:       d,
	}
}

// GetUser 获取管理员个人信息
func (ac *AdminController) GetUser(c *gin.Context) {
	id, err := GetUserID(c)
	if err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "请求无效"))
		return
	}
	response := ac.userService.GetUserByID(id)
	c.JSON(http.StatusOK, response)
}

// GetUsers 获取所有用户信息
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
	logrus.Debug(payload)
	if err != nil {
		request.UserID = 0
	} else {
		request.UserID = payload.ID
	}

	response := ac.videoSearchService.GetVideos(&request)
	c.JSON(http.StatusOK, response)
}

// UpdateVideo 更新视频信息
func (ac *AdminController) UpdateVideo(c *gin.Context) {
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
	if err := c.ShouldBind(&request); err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "请求参数错误"))
		return
	}

	response := ac.userVideoService.DeleteVideoByAdmin(&request)
	c.JSON(http.StatusOK, response)
}

// GetRealTimeData 获取实时数据(在线人数, cpu, 内存)
func (ac *AdminController) GetRealTimeData(c *gin.Context) {
	resonse := ac.statsService.GetRealTimeData()
	c.JSON(http.StatusOK, resonse)
}

// GetHistoricalData 获取历史数据(新增用户数, 新增视频浏览量, 每天登录用户数)
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
