package controller

import (
	"fmt"
	"net/http"
	"videohub/internal/service"

	"github.com/gin-gonic/gin"
)

// AdminController 包含用户相关的服务
type AdminController struct {
	userAvatarService *service.UserAvatar // 用户头像服务
	userListService   *service.UserList   // 用户列表服务
	userService       *service.User       // 用户服务
}

// NewAdminController 创建一个新的 AdminController 实例
func NewAdminController(uas *service.UserAvatar, uls *service.UserList, us *service.User) *AdminController {
	return &(AdminController{userAvatarService: uas, userListService: uls, userService: us})
}

// GetUser 获取管理员个人信息
func (ac *AdminController) GetUser(c *gin.Context) {
	id, err := GetUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果获取 ID 失败，返回 HTTP 400
		return
	}
	fmt.Println("Get videos for user", id)
	// TODO
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
