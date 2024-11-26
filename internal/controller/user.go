package controller

import (
	"errors"
	"fmt"
	"net/http"
	"videohub/global"
	"videohub/internal/service"
	"videohub/internal/utils/user"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userAvatarService *service.UserAvatar // 用户头像服务
	userListService   *service.UserList   // 用户列表服务
	userService       *service.User       // 用户服务
}

func NewUserController(uas *service.UserAvatar, uls *service.UserList, us *service.User) *UserController {
	return &(UserController{userAvatarService: uas, userListService: uls, userService: us})
}

func GetUserID(c *gin.Context) (uint, error) {
	idValue, exists := c.Get("id")
	if !exists {
		return 0, errors.New("上下文中不存在用户 ID")
	}
	id, ok := idValue.(uint)
	if !ok {
		return 0, errors.New("用户 ID 类型错误")
	}
	return id, nil
}

// Login 用户登录处理函数，返回范文令牌和刷新令牌给前端
func (uc *UserController) Login(c *gin.Context) {
	var request user.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求的JSON格式无效或缺少必需字段"})
		return
	}
	response := uc.userService.Login(&request)
	c.JSON(response.StatusCode, response)
}

// AccessToken 获取访问令牌处理函数，返回访问令牌给前端
func (uc *UserController) AccessToken(c *gin.Context) {
	var request user.AccessTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求的JSON格式无效或缺少必需字段"})
		return
	}
	response := uc.userService.AccessToken(&request)
	c.JSON(response.StatusCode, response)
}

// GetUserInfo 获取某个用户信息 (根据用户的 access_token 拿到 id)
func (uc *UserController) GetUser(c *gin.Context) {
	id, err := GetUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // 如果获取 ID 失败，返回 HTTP 500
		return
	}

	response := uc.userService.GetUserByID(id)
	c.JSON(response.StatusCode, response)
}

// CreateUser 创建新用户
func (uc *UserController) CreateUser(c *gin.Context) {
	var request user.CreateUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求的JSON格式无效或缺少必需字段"})
		return
	}

	if err := global.Validate.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱格式无效"}) // 如果邮箱格式无效，返回 HTTP 400
		return
	}

	response := uc.userService.CreateUser(&request)
	c.JSON(response.StatusCode, response)
}

// UpdateUser 根据用户 ID 更新用户信息
func (uc *UserController) UpdateUser(c *gin.Context) {
	id, err := GetUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果获取 ID 失败，返回 HTTP 400
		return
	}

	var request user.UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求的JSON格式无效或缺少必需字段"}) // 如果解析 JSON 失败，返回 HTTP 400
		return
	}

	var fileds []string
	if request.Email != "" {
		if err := global.Validate.Struct(request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱格式无效"}) // 如果邮箱格式无效，返回 HTTP 400
			return
		}
		fileds = append(fileds, "email")
	}

	if request.Username != "" {
		fileds = append(fileds, "username")
	}

	response := uc.userService.UpdateUser(id, fileds, &request)
	c.JSON(response.StatusCode, response)
}

// DeleteUser 根据用户 ID 删除用户
func (uc *UserController) DeleteUser(c *gin.Context) {
	id, err := GetUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果获取 ID 失败，返回 HTTP 400
		return
	}

	response := uc.userService.DeleteUser(id)
	c.JSON(response.StatusCode, response)
}

// UploadAvatar 上传用户头像（携带头像数据）
func (uc *UserController) UploadAvatar(c *gin.Context) {
	id, err := GetUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果获取 ID 失败，返回 HTTP 400
		return
	}

	var request user.UploadAvatarRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求的Form格式无效或缺少必需字段"})
		return
	}

	response := uc.userAvatarService.UploadUserAvatar(id, &request)
	c.JSON(response.StatusCode, response)
}

// UpdatePassword 修改用户密码
func (uc *UserController) UpdatePassword(c *gin.Context) {
	id, err := GetUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果获取 ID 失败，返回 HTTP 400
		return
	}

	var request user.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求的JSON格式无效或缺少必需字段"})
		return
	}

	if err := global.Validate.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱格式无效"}) // 如果邮箱格式无效，返回 HTTP 400
		return
	}

	response := uc.userService.UpdateUserPassword(id, &request)
	c.JSON(response.StatusCode, response)
}

// SendEmailVerification 发送验证码到邮箱
// 采用 gomail + SMTP 发送邮件 (163 邮箱)
func (uc *UserController) SendEmailVerification(c *gin.Context) {
	var request user.SendEmailVerificationRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求的JSON格式无效或缺少必需字段"}) // 如果解析 JSON 失败，返回 HTTP 400
		return
	}

	if err := global.Validate.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱格式无效"}) // 如果邮箱格式无效，返回 HTTP 400
		return
	}

	response := uc.userService.SendEmailVerification(&request)
	c.JSON(response.StatusCode, response)
}

// GetVideos 获取用户上传的视频列表
func (uc *UserController) GetVideos(c *gin.Context) {
	id, err := GetUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果获取 ID 失败，返回 HTTP 400
		return
	}
	fmt.Println("Get videos for user", id)
	// TODO
}

// DeleteVideo 删除用户上传的视频
func (uc *UserController) DeleteVideo(c *gin.Context) {
	id, err := GetUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果获取 ID 失败，返回 HTTP 400
		return
	}
	fmt.Println("Get videos for user", id)
	// TODO
}

// GetCollections 获取用户收藏的视频列表
func (uc *UserController) GetCollections(c *gin.Context) {
	id, err := GetUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果获取 ID 失败，返回 HTTP 400
		return
	}
	fmt.Println("Get videos for user", id)
	// TODO
}

// UpdateCollections 更新用户收藏的视频
func (uc *UserController) UpdateCollections(c *gin.Context) {
	id, err := GetUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果获取 ID 失败，返回 HTTP 400
		return
	}
	fmt.Println("Get videos for user", id)
	// TODO
}

// DeleteCollections 删除用户收藏的视频
func (uc *UserController) DeleteCollections(c *gin.Context) {
	id, err := GetUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果获取 ID 失败，返回 HTTP 400
		return
	}
	fmt.Println("Get videos for user", id)
	// TODO
}
