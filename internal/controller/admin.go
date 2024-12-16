package controller

import (
	"fmt"
	"net/http"
	"videohub/internal/service"
	"videohub/internal/utils"
	"videohub/internal/utils/admin"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AdminController 包含用户相关的服务
type AdminController struct {
	userAvatarService *service.UserAvatar // 用户头像服务
	userListService   *service.UserList   // 用户列表服务
	userService       *service.User       // 用户服务
	dataService       *service.Stats      // 数据服务
}

// NewAdminController 创建一个新的 AdminController 实例
func NewAdminController(uas *service.UserAvatar, uls *service.UserList, us *service.User, d *service.Stats) *AdminController {
	return &(AdminController{userAvatarService: uas, userListService: uls, userService: us, dataService: d})
}

// GetUser 获取管理员个人信息
func (ac *AdminController) GetUser(c *gin.Context) {
	id, err := GetUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "请求无效"))
		return
	}
	fmt.Println("Get videos for user", id)
	// TODO
}

// GetUsers 获取用户信息
func (ac *AdminController) GetUsers(c *gin.Context) {
	// TODO
}

// CreateUser 创建用户
func (ac *AdminController) CreateUser(c *gin.Context) {
	// TODO
}

// UpdateUser 更新用户信息
func (ac *AdminController) UpdateUser(c *gin.Context) {
	// TODO
}

// GetVideos 获取视频列表
func (ac *AdminController) GetVideos(c *gin.Context) {
	// TODO
}

// UpdateVideo 更新视频信息
func (ac *AdminController) UpdateVideo(c *gin.Context) {
	// TODO
}

// DeleteVideo 删除视频
func (ac *AdminController) DeleteVideo(c *gin.Context) {
	// TODO
}

func (ac *AdminController) GetRealTimeData(c *gin.Context) {
	resonse := ac.dataService.GetRealTimeData()
	c.JSON(http.StatusOK, resonse)
}

func (ac *AdminController) GetHistoricalData(c *gin.Context) {
	var request admin.GetHistoricalDataRequest
	if err := c.ShouldBind(&request); err != nil {
		logrus.Debug(err.Error())
		c.JSON(http.StatusOK, utils.Error(http.StatusBadRequest, "请求无效"))
		return
	}
	response := ac.dataService.GetHistoricalData(&request)
	c.JSON(http.StatusOK, response)
}
