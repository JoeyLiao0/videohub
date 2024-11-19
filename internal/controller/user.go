package controller

import (
	"fmt"
	"net/http"
	"videohub/internal/model"
	"videohub/internal/service"
	"videohub/internal/utils"

	"github.com/gin-gonic/gin"
)

// UserController handles user-related operations and services.
// It includes services for managing user avatars, user lists, and general user operations.
type UserController struct {
	userAvatarService *service.UserAvatar // 用户头像服务
	userListService   *service.UserList   // 用户列表服务
	userService       *service.User       // 用户服务
}

// NewUserController creates a new instance of UserController with the provided services.
// Parameters:
//   - uas: a pointer to UserAvatar service
//   - uls: a pointer to UserList service
//   - us: a pointer to User service
//
// Returns:
//   - A pointer to a UserController instance initialized with the provided services.
func NewUserController(uas *service.UserAvatar, uls *service.UserList, us *service.User) *UserController {
	return &(UserController{userAvatarService: uas, userListService: uls, userService: us})
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
func (uc *UserController) getUserID(c *gin.Context) (uint64, error) {
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

// Login 用户登录处理函数，返回范文令牌和刷新令牌给前端
func (uc *UserController) Login(c *gin.Context) {
	uc.userService.Login(c)
}

// AccessToken 获取访问令牌处理函数，返回访问令牌给前端
func (uc *UserController) AccessToken(c *gin.Context) {
	uc.userService.AccessToken(c)
}

// GetUserInfo 获取某个用户信息 (根据用户的 access_token 拿到 id)
func (uc *UserController) GetUser(c *gin.Context) {
	id, err := uc.getUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // 如果获取 ID 失败，返回 HTTP 500
		return
	}

	user, err := uc.userService.GetUserByID(id) // 调用服务层根据 ID 获取用户信息
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()}) // 404
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user}) // 成功则返回 HTTP 200 和用户信息
}

// CreateUser 创建新用户
func (uc *UserController) CreateUser(c *gin.Context) {
	var inputs service.CreateUserInput
	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"}) // 如果解析 JSON 失败，返回 HTTP 400
		return
	}

	var newUser model.User
	if err := uc.userService.CreateUser(inputs, &newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果创建用户失败，返回 HTTP 400 和错误信息
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": newUser}) // 成功则返回 HTTP 201 和用户信息
}

// UpdateUser 根据用户 ID 更新用户信息
func (uc *UserController) UpdateUser(c *gin.Context) {
	var request service.UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"}) // 如果解析 JSON 失败，返回 HTTP 400
		return
	}

	id, err := uc.getUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果获取 ID 失败，返回 HTTP 400
		return
	}

	oldUser, err := uc.userService.GetUserByID(uint64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()}) // 如果用户不存在，返回 HTTP 404
		return
	}

	// TODO: 补充更新字段的逻辑
	if request.Username != "" && request.Username != oldUser.Username {
		oldUser.Username = request.Username
	}

	if err := uc.userService.UpdateUser(oldUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果更新用户失败，返回 HTTP 400 和错误信息
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": oldUser}) // 成功则返回 HTTP 200 和更新后的用户信息
}

// DeleteUser 根据用户 ID 删除用户
func (uc *UserController) DeleteUser(c *gin.Context) {
	id, err := uc.getUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果获取 ID 失败，返回 HTTP 400
		return
	}

	if err := uc.userService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果删除用户失败，返回 HTTP 400 和错误信息
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"}) // 成功则返回 HTTP 200 和删除成功消息
}

// UploadAvatar 上传用户头像（携带头像数据）
func (uc *UserController) UploadAvatar(c *gin.Context) {
	id, err := uc.getUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果获取 ID 失败，返回 HTTP 400
		return
	}

	// 获取上传的头像文件 (multipart form-data)
	// TODO: 文件类型、大小检查
	file, _, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload error"})
		return
	}

	// 调用服务层处理头像上传
	if err := uc.userAvatarService.UploadUserAvatar(uint(id), file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Avatar uploaded successfully"})
}

// UpdatePassword 修改用户密码
func (uc *UserController) UpdatePassword(c *gin.Context) {
	id, err := uc.getUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果获取 ID 失败，返回 HTTP 400
		return
	}

	var passwordData struct {
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&passwordData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := uc.userService.UpdateUserPassword(uint64(id), passwordData.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// SendEmailVerification 发送验证码到邮箱
// TODO: 采用 gomail + SMTP 发送邮件 (qq 邮箱)
func (uc *UserController) SendEmailVerification(c *gin.Context) {
	var emailData struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&emailData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := utils.SendEmailVerification(emailData.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent"})
}

// GetVideos 获取用户上传的视频列表
func (uc *UserController) GetVideos(c *gin.Context) {
	id, err := uc.getUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果获取 ID 失败，返回 HTTP 400
		return
	}
	fmt.Println("Get videos for user", id)
	// TODO
}

// DeleteVideos 删除用户上传的视频
func (uc *UserController) DeleteVideos(c *gin.Context) {
	id, err := uc.getUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果获取 ID 失败，返回 HTTP 400
		return
	}
	fmt.Println("Get videos for user", id)
	// TODO
}

// GetCollections 获取用户收藏的视频列表
func (uc *UserController) GetCollections(c *gin.Context) {
	id, err := uc.getUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果获取 ID 失败，返回 HTTP 400
		return
	}
	fmt.Println("Get videos for user", id)
	// TODO
}

// UpdateCollections 更新用户收藏的视频
func (uc *UserController) UpdateCollections(c *gin.Context) {
	id, err := uc.getUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果获取 ID 失败，返回 HTTP 400
		return
	}
	fmt.Println("Get videos for user", id)
	// TODO
}

// DeleteCollections 删除用户收藏的视频
func (uc *UserController) DeleteCollections(c *gin.Context) {
	id, err := uc.getUserID(c) // 从上下文中获取用户 ID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果获取 ID 失败，返回 HTTP 400
		return
	}
	fmt.Println("Get videos for user", id)
	// TODO
}
