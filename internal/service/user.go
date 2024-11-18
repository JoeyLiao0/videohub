package service

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
	"videohub/config"
	"videohub/internal/model"
	"videohub/internal/repository"
	"videohub/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	//用户服务，用到user表、collection表、video表操作
	userRepo       *repository.User
	collectionRepo *repository.Collection
	videoRepo      *repository.Video
}

type CreateUserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Avatar   string `json:"avatar"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Status   uint8  `json:"status"`
	Role     uint8  `json:"role"`
}

type UpdateUserResponse struct {
	
}

// 工厂函数，返回单例的服务层操作对象
func NewUser(ur *repository.User, cr *repository.Collection, vr *repository.Video) *User {
	return &(User{userRepo: ur, collectionRepo: cr, videoRepo: vr})
}

// 服务函数追加在下面
func (us *User) Login(c *gin.Context) {
	// 或者使用 map[string]interface{}
	var inputs struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 400
		return
	}

	var user model.User
	if err := us.userRepo.SearchByUsername(inputs.Username, &user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong credentials"}) // 401
		return
	}
	log.Println(user.Password, inputs.Password)
	if user.Password != utils.HashPassword(inputs.Password, user.Salt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong credentials"}) // 401
		return
	}

	accessToken, err := utils.GenerateJWT(utils.Payload{ID: user.ID, Role: user.Role}, config.AppConfig.JWT.AccessTokenSecret, config.AppConfig.JWT.AccessTokenExpire)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"}) // 500
		return
	}

	refreshToken, err := utils.GenerateJWT(utils.Payload{ID: user.ID, Role: user.Role}, config.AppConfig.JWT.RefreshTokenSecret, config.AppConfig.JWT.RefreshTokenExpire)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"}) // 500
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (us *User) AccessToken(c *gin.Context) {
	var inputs struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 400
		return
	}

	payload, err := utils.ParseJWT(inputs.RefreshToken, config.AppConfig.JWT.RefreshTokenSecret)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()}) // 401
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"}) // 500
		}
		return
	}

	var user model.User
	if err := us.userRepo.SearchById(payload.ID, &user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()}) // 401
		return
	}

	accessToken, err := utils.GenerateJWT(payload, config.AppConfig.JWT.AccessTokenSecret, config.AppConfig.JWT.AccessTokenExpire)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"}) // 500
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

// GetUserByID 根据用户 ID 获取单个用户信息
func (us *User) GetUserByID(userID uint64) (*model.User, error) {
	// 调用 repository 查询用户
	user, err := us.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser 创建新用户
// TODO: 是否创建钩子函数
func (us *User) CreateUser(inputs CreateUserInput, newUser *model.User) error {
	// 1. 检查用户名是否唯一
	existingUser, err := us.userRepo.FindUserByUsername(inputs.Username)
	if err != nil {
		return fmt.Errorf("查询用户名时发生错误: %v", err)
	}
	if existingUser != nil {
		return fmt.Errorf("用户名已存在")
	}

	// 2. 检查邮箱是否唯一
	existingEmailUser, err := us.userRepo.FindUserByEmail(inputs.Email)
	if err != nil {
		return fmt.Errorf("查询邮箱时发生错误: %v", err)
	}
	if existingEmailUser != nil {
		return fmt.Errorf("邮箱已被注册")
	}

	// 3. 创建新用户
	newUser.Username = inputs.Username
	newUser.Salt = utils.GenerateSalt(16)
	newUser.Password = utils.HashPassword(inputs.Password, newUser.Salt)
	newUser.CreatedAt = time.Now().UnixMilli()
	newUser.Avatar = inputs.Avatar
	newUser.Email = inputs.Email
	newUser.Status = 1
	newUser.Role = 0

	// 4. 存储用户信息
	if err := us.userRepo.CreateUser(newUser); err != nil {
		return err
	}

	return nil
}

// UpdateUser 更新用户信息
func (us *User) UpdateUser(updatedUser *model.User) error {
	// 调用 repository 更新用户信息
	if err := us.userRepo.UpdateUser(updatedUser); err != nil {
		return err
	}
	return nil
}

// DeleteUser 根据用户 ID 删除用户
func (us *User) DeleteUser(userID uint) error {
	// 调用 repository 删除用户
	if err := us.userRepo.DeleteUser(userID); err != nil {
		return err
	}
	return nil
}

// UpdateUserPassword 更新用户密码
func (us *User) UpdateUserPassword(userID uint64, newPassword string) error {
	// 查找用户
	user, err := us.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	// 生成新的盐值并加密新密码
	user.Salt = utils.GenerateSalt(16)
	user.Password = utils.HashPassword(newPassword, user.Salt)

	// 更新用户密码
	if err := us.userRepo.UpdateUser(user); err != nil {
		return err
	}
	return nil
}
