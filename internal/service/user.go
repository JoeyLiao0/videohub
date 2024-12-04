package service

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"videohub/config"
	"videohub/global"
	"videohub/internal/model"
	"videohub/internal/repository"
	"videohub/internal/utils"
	"videohub/internal/utils/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
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

func (us *User) Login(request *user.LoginRequest) *utils.Response {
	var result model.User
	if err := us.userRepo.Search(map[string]interface{}{"email": request.Email}, 1, &result); err != nil {
		logrus.Debug(err.Error())
		return utils.Error(http.StatusUnauthorized, "邮箱未注册")
	}

	if result.Password != utils.HashPassword(request.Password, result.Salt) {
		logrus.Debug("password error")
		return utils.Error(http.StatusBadRequest, "密码错误")
	}

	if result.Status == 2 {
		logrus.Debug("user is banned")
		return utils.Error(http.StatusUnauthorized, "用户已注销")
	}

	accessToken, err := utils.GenerateJWT(utils.Payload{ID: result.ID, Role: result.Role}, config.AppConfig.JWT.AccessTokenSecret, config.AppConfig.JWT.AccessTokenExpire)
	if err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	refreshToken, err := utils.GenerateJWT(utils.Payload{ID: result.ID, Role: result.Role}, config.AppConfig.JWT.RefreshTokenSecret, config.AppConfig.JWT.RefreshTokenExpire)
	if err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	logrus.Debug("Login successfully")
	return utils.Ok(http.StatusOK, user.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (us *User) AccessToken(request *user.AccessTokenRequest) *utils.Response {
	payload, err := utils.ParseJWT(request.RefreshToken, config.AppConfig.JWT.RefreshTokenSecret)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			logrus.Debug(err.Error())
			return utils.Error(http.StatusBadRequest, "令牌已过期")
		} else {
			logrus.Error(err.Error())
			return utils.Error(http.StatusInternalServerError, "服务器内部错误")
		}
	}

	if count, err := us.userRepo.Count(map[string]interface{}{"id": payload.ID}); err != nil || count == 0 {
		logrus.Debug("user not found")
		return utils.Error(http.StatusUnauthorized, "未授权")
	}

	accessToken, err := utils.GenerateJWT(payload, config.AppConfig.JWT.AccessTokenSecret, config.AppConfig.JWT.AccessTokenExpire)
	if err != nil {
		logrus.Debug(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	logrus.Debug("Access token refreshed successfully")
	return utils.Ok(http.StatusOK, user.AccessTokenResponse{AccessToken: accessToken})
}

// GetUserByID 根据用户 ID 获取单个用户信息
func (us *User) GetUserByID(id uint) *utils.Response {
	var response user.GetUserResponse
	if err := us.userRepo.Search(map[string]interface{}{"id": id}, 1, &response.User); err != nil {
		logrus.Debug(err.Error())
		return utils.Error(http.StatusInternalServerError, err.Error())
	}

	logrus.Debug("User retrieved successfully")
	return utils.Ok(http.StatusOK, response)
}

// CreateUser 创建新用户
func (us *User) CreateUser(request *user.CreateUserRequest) *utils.Response {
	var username string
	for {
		username, _ = utils.GenerateUsername(12)
		if count, err := us.userRepo.Count(map[string]interface{}{"username": username}); err != nil || count != 0 {
			logrus.Debug("username exists")
		} else {
			break
		}
	}

	if count, err := us.userRepo.Count(map[string]interface{}{"email": request.Email}); err != nil || count != 0 {
		logrus.Debug("email exists")
		return utils.Error(http.StatusBadRequest, "邮箱已被注册")
	}

	if err := utils.VerifyEmailCode(request.Email, request.Code); err != nil {
		logrus.Debug(err.Error())
		return utils.Error(http.StatusBadRequest, "验证码错误")
	}

	var newUser model.User
	newUser.Username = username
	newUser.Salt = utils.GenerateSalt(16)
	newUser.Password = utils.HashPassword(request.Password, newUser.Salt)
	newUser.Email = request.Email
	newUser.Avatar = fmt.Sprintf("%s/%s/%s", config.AppConfig.Storage.Base, config.AppConfig.Storage.Images, "tourist.jpeg")

	if err := us.userRepo.Create(&newUser); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	logrus.Debugf("User %s created successfully", newUser.Username)
	return utils.Success(http.StatusOK)
}

// UpdateUser 更新用户信息
func (us *User) UpdateUser(id uint, fileds interface{}, request *user.UpdateUserRequest) *utils.Response {
	if request.Email != "" {
		if count, err := us.userRepo.Count(map[string]interface{}{"email": request.Email}); err != nil || count != 0 {
			logrus.Debug("email exists")
			return utils.Error(http.StatusBadRequest, "邮箱已被注册")
		}

		if request.Code == "" {
			logrus.Debug("code is empty")
			return utils.Error(http.StatusBadRequest, "验证码为空")
		}

		if err := utils.VerifyEmailCode(request.Email, request.Code); err != nil {
			logrus.Debug(err.Error())
			return utils.Error(http.StatusBadRequest, "验证码错误")
		}
	}

	if request.Username != "" {
		if count, err := us.userRepo.Count(map[string]interface{}{"username": request.Username}); err != nil || count != 0 {
			logrus.Debug("username exists")
			return utils.Error(http.StatusBadRequest, "用户名已被注册")
		}
	}

	if err := us.userRepo.Update(map[string]interface{}{"id": id}, fileds, request); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	logrus.Debug("User updated successfully")
	return utils.Success(http.StatusOK)
}

// DeleteUser 根据用户 ID 删除用户
func (us *User) DeleteUser(id uint) *utils.Response {
	values := map[string]interface{}{
		"status": 2,
		"avatar": fmt.Sprintf("%s/%s/%s",config.AppConfig.Storage.Base, config.AppConfig.Storage.Images, "logout.jpeg"),
	}
	fileds := []string{"status", "avatar"}
	if err := us.userRepo.Update(map[string]interface{}{"id": id}, fileds, values); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	logrus.Debug("User deleted successfully")
	return utils.Success(http.StatusOK)
}

// UpdateUserPassword 更新用户密码
func (us *User) UpdateUserPassword(id uint, request *user.UpdatePasswordRequest) *utils.Response {
	var result struct {
		Email    string
		Salt     string
		Password string
	}
	if err := us.userRepo.Search(map[string]interface{}{"id": id}, 1, &result); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusUnauthorized, "未授权")
	}

	if utils.HashPassword(request.Password, result.Salt) != result.Password {
		logrus.Debug("wrong password")
		return utils.Error(http.StatusBadRequest, "密码错误")
	}

	if err := utils.VerifyEmailCode(result.Email, request.Code); err != nil {
		logrus.Debug(err.Error())
		return utils.Error(http.StatusBadRequest, "验证码错误")
	}

	salt := utils.GenerateSalt(16)
	password := utils.HashPassword(request.NewPassword, salt)

	if err := us.userRepo.Update(map[string]interface{}{"id": id}, []string{"salt", "password"}, map[string]interface{}{"salt": salt, "password": password}); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	logrus.Debug("Password updated successfully")
	return utils.Success(http.StatusOK)
}

func (us *User) SendEmailVerification(request *user.SendEmailVerificationRequest) *utils.Response {
	code := utils.GenerateCode(6)
	global.Rdb.Set(global.Ctx, request.Email, code, time.Minute*time.Duration(config.AppConfig.Email.Expiration))
	if err := utils.SendEmailVerification(request.Email, code); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	logrus.Debug("Email verification sent successfully")
	return utils.Success(http.StatusOK)
}
