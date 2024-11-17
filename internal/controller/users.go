package controller

import (
	"net/http"
	"strconv"
	"videohub/internal/model"
	"videohub/internal/service"
	"videohub/internal/utils"

	"github.com/gin-gonic/gin"
)

// Users 控制器结构体，包含用户相关的服务
type Users struct {
	userAvatarService *service.User_avatar // 用户头像服务
	userListService   *service.User_list   // 用户列表服务
	userService       *service.User        // 用户服务
}

// NewUsers 创建一个新的 Users 控制器实例
func NewUsers(uas *service.User_avatar, uls *service.User_list, us *service.User) *Users {
	return &(Users{userAvatarService: uas, userListService: uls, userService: us})
}

// Login 用户登录处理函数，返回范文令牌和刷新令牌给前端
func (uc *Users) Login(c *gin.Context) {
	uc.userService.Login(c)
}

// AccessToken 获取访问令牌处理函数，返回访问令牌给前端
func (uc *Users) AccessToken(c *gin.Context) {
	uc.userService.AccessToken(c)
}

// GetAllUsers 获取所有普通用户的信息
func (uc *Users) GetAllUsers(c *gin.Context) {
	users, err := uc.userListService.GetAllUsers() // 调用服务层获取所有用户信息
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // 如果有错误，返回 HTTP 500 和错误信息
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users}) // 成功则返回 HTTP 200 和用户列表
}

// CreateUser 创建新用户
func (uc *Users) CreateUser(c *gin.Context) {
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

// GetUserByID 根据用户 ID 获取用户信息
func (uc *Users) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")                      // 从 URL 参数中获取用户 ID 字符串
	id, err := strconv.ParseUint(idStr, 10, 32) // 将字符串转换为 uint 类型
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"}) // 如果转换失败，返回 HTTP 400
		return
	}

	user, err := uc.userService.GetUserByID(uint64(id)) // 调用服务层根据 ID 获取用户信息
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()}) // 如果用户不存在，返回 HTTP 404
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user}) // 成功则返回 HTTP 200 和用户信息
}

// UpdateUser 根据用户 ID 更新用户信息
func (uc *Users) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")                      // 从 URL 参数中获取用户 ID 字符串
	id, err := strconv.ParseUint(idStr, 10, 32) // 将字符串转换为 uint 类型
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"}) // 如果转换失败，返回 HTTP 400
		return
	}

	var inputs service.UpdateUserInput
	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"}) // 如果解析 JSON 失败，返回 HTTP 400
		return
	}

	oldUser, err := uc.userService.GetUserByID(uint64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()}) // 如果用户不存在，返回 HTTP 404
		return
	}

	roleValue, exists := c.Get("role")
	role, ok := roleValue.(uint8)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid role type"})
		return
	}
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Role not found"})
		return
	}

	if inputs.Username != "" && inputs.Username != oldUser.Username {
		oldUser.Username = inputs.Username
	}

	if role == 0 { // 普通用户
		if inputs.Email != "" && inputs.Email != oldUser.Email {
			oldUser.Email = inputs.Email
			// TODO: 发送验证邮件
		}
	} else if role == 1 { // 管理员
		if inputs.Email != "" && inputs.Email != oldUser.Email {
			oldUser.Email = inputs.Email
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := uc.userService.UpdateUser(oldUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果更新用户失败，返回 HTTP 400 和错误信息
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": oldUser}) // 成功则返回 HTTP 200 和更新后的用户信息
}

// DeleteUser 根据用户 ID 删除用户
func (uc *Users) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")                      // 从 URL 参数中获取用户 ID 字符串
	id, err := strconv.ParseUint(idStr, 10, 32) // 将字符串转换为 uint 类型
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"}) // 如果转换失败，返回 HTTP 400
		return
	}

	roleValue, exists := c.Get("role")
	role, ok := roleValue.(uint8)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid role type"})
		return
	}
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Role not found"})
		return
	}

	// 只有管理员可以删除用户
	if role != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := uc.userService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 如果删除用户失败，返回 HTTP 400 和错误信息
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"}) // 成功则返回 HTTP 200 和删除成功消息
}

// UploadAvatar 上传用户头像（携带头像数据）
func (uc *Users) UploadAvatar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 获取上传的头像文件 (multipart form-data)
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
func (uc *Users) UpdatePassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
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
func (uc *Users) SendEmailVerification(c *gin.Context) {
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
