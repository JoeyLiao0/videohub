package service

import (
	"errors"
	"log"
	"net/http"
	"time"
	"videohub/config"
	"videohub/internal/model"
	"videohub/internal/repository"
	"videohub/internal/utils"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	//用户服务，用到user表、collection表、video表操作
	userRepo       *repository.User
	collectionRepo *repository.Collection
	videoRepo      *repository.Video
}

// 工厂函数，返回单例的服务层操作对象
func NewUser(ur *repository.User, cr *repository.Collection, vr *repository.Video) *User {
	return &(User{userRepo: ur, collectionRepo: cr, videoRepo: vr})
}

func (us *User) Login(request utils.LoginRequest) *utils.Response {
	var user model.User
	if err := us.userRepo.SearchByEmail(request.Email, &user); err != nil {
		return utils.Error(http.StatusUnauthorized, "邮箱未注册")
	}
	log.Println(user.Password, request.Password)
	if user.Password != utils.HashPassword(request.Password, user.Salt) {
		return utils.Error(http.StatusUnauthorized, "邮箱或密码错误")
	}

	accessToken, err := utils.GenerateJWT(utils.Payload{ID: user.ID, Role: user.Role}, config.AppConfig.JWT.AccessTokenSecret, config.AppConfig.JWT.AccessTokenExpire)
	if err != nil {
		return utils.Error(http.StatusInternalServerError, "生成访问令牌失败")
	}

	refreshToken, err := utils.GenerateJWT(utils.Payload{ID: user.ID, Role: user.Role}, config.AppConfig.JWT.RefreshTokenSecret, config.AppConfig.JWT.RefreshTokenExpire)
	if err != nil {
		return utils.Error(http.StatusInternalServerError, "生成刷新令牌失败")
	}

	return utils.Ok(http.StatusOK, utils.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (us *User) AccessToken(request utils.AccessTokenRequest) *utils.Response {
	payload, err := utils.ParseJWT(request.RefreshToken, config.AppConfig.JWT.RefreshTokenSecret)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return utils.Error(http.StatusUnauthorized, "刷新令牌已过期")
		} else {
			return utils.Error(http.StatusInternalServerError, "解析刷新令牌失败")
		}
	}

	var user model.User
	if err := us.userRepo.SearchById(payload.ID, &user); err != nil {
		return utils.Error(http.StatusUnauthorized, "用户不存在")
	}

	accessToken, err := utils.GenerateJWT(payload, config.AppConfig.JWT.AccessTokenSecret, config.AppConfig.JWT.AccessTokenExpire)
	if err != nil {
		return utils.Error(http.StatusInternalServerError, "生成访问令牌失败")
	}

	return utils.Ok(http.StatusOK, utils.AccessTokenResponse{AccessToken: accessToken})
}

// GetUserByID 根据用户 ID 获取单个用户信息
func (us *User) GetUserByID(userID uint64) *utils.Response {
	// 调用 repository 查询用户
	user, err := us.userRepo.GetUserByID(userID)
	if err != nil {
		return utils.Error(http.StatusNotFound, err.Error())
	}
	return utils.Ok(http.StatusOK, user)
}

// CreateUser 创建新用户
// TODO: 是否创建钩子函数
func (us *User) CreateUser(request utils.CreateUserRequest, newUser *model.User) *utils.Response {
	// 1. 检查用户名是否唯一
	existingUser, err := us.userRepo.FindUserByUsername(request.Username)
	if err != nil {
		return utils.Error(http.StatusInternalServerError, "查询用户名时发生错误")
	}
	if existingUser != nil {
		return utils.Error(http.StatusBadRequest, "用户名已被注册")
	}

	// 2. 检查邮箱是否唯一
	existingEmailUser, err := us.userRepo.FindUserByEmail(request.Email)
	if err != nil {
		return utils.Error(http.StatusInternalServerError, "查询邮箱时发生错误")
	}
	if existingEmailUser != nil {
		return utils.Error(http.StatusBadRequest, "邮箱已被注册")
	}

	// 3. 创建新用户
	newUser.Username = request.Username
	newUser.Salt = utils.GenerateSalt(16)
	newUser.Password = utils.HashPassword(request.Password, newUser.Salt)
	newUser.CreatedAt = time.Now().UnixMilli()
	newUser.Avatar = request.Avatar
	newUser.Email = request.Email
	newUser.Status = 1
	newUser.Role = 0

	// 4. 存储用户信息
	if err := us.userRepo.CreateUser(newUser); err != nil {
		return utils.Error(http.StatusInternalServerError, "创建用户时发生错误")
	}

	return utils.Ok(http.StatusCreated, "用户创建成功")
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
