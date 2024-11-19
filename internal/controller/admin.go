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

// getUserID retrieves the user ID from the given gin.Context.
// It returns the user ID as a uint64 and an error if the ID is not found
// in the context or if the ID is of an invalid type.
//
// Parameters:
//   - c: The gin.Context from which to retrieve the user ID.
//
// Returns:
//   - uint64: The user ID.
//   - error: An error if the user ID is not found or is of an invalid type.
func (ac *AdminController) getUserID(c *gin.Context) (uint64, error) {
	idValue, exists := c.Get("id")
	if !exists {
		return 0, fmt.Errorf("User ID not found in context")
	}
	id, ok := idValue.(uint64)
	if !ok {
		return 0, fmt.Errorf("Invalid user ID type")
	}
	return id, nil
}

// GetUser 获取管理员个人信息
func (ac *AdminController) GetUser(c *gin.Context) {
	id, err := ac.getUserID(c) // 从上下文中获取用户 ID
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

// CreateUsers 创建用户
func (ac *AdminController) CreateUsers(c *gin.Context) {
	// TODO
}

// UpdateUsers 更新用户信息
func (ac *AdminController) UpdateUsers(c *gin.Context) {
	// TODO
}

// GetVideos 获取视频列表
func (ac *AdminController) GetVideos(c *gin.Context) {
	// TODO
}

// UpdateVideos 更新视频信息
func (ac *AdminController) UpdateVideos(c *gin.Context) {
	// TODO
}

// DeleteVideos 删除视频
func (ac *AdminController) DeleteVideos(c *gin.Context) {
	// TODO
}
