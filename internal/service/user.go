package service

import (
	"errors"
	"log"
	"net/http"
	"time"
	"videohub/config"
	"videohub/global"
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
	if err := us.userRepo.Search(map[string]interface{}{"email": request.Email}, 1, &user); err != nil {
		return utils.Error(http.StatusUnauthorized, "邮箱未注册")
	}

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

	if count, err := us.userRepo.Count(map[string]interface{}{"id": payload.ID}); err != nil || count == 0 {
		return utils.Error(http.StatusUnauthorized, "用户不存在")
	}

	accessToken, err := utils.GenerateJWT(payload, config.AppConfig.JWT.AccessTokenSecret, config.AppConfig.JWT.AccessTokenExpire)
	if err != nil {
		return utils.Error(http.StatusInternalServerError, "生成访问令牌失败")
	}

	return utils.Ok(http.StatusOK, utils.AccessTokenResponse{AccessToken: accessToken})
}

// GetUserByID 根据用户 ID 获取单个用户信息
func (us *User) GetUserByID(id uint64) *utils.Response {
	var response utils.GetUserResponse
	if err := us.userRepo.Search(map[string]interface{}{"id": id}, 1, &response); err != nil {
		return utils.Error(http.StatusNotFound, err.Error())
	}
	return utils.Ok(http.StatusOK, response)
}

// CreateUser 创建新用户
func (us *User) CreateUser(request utils.CreateUserRequest) *utils.Response {
	if count, err := us.userRepo.Count(map[string]interface{}{"username": request.Username}); err != nil || count != 0 {
		log.Println(request)
		return utils.Error(http.StatusBadRequest, "该用户名已存在")
	}

	if count, err := us.userRepo.Count(map[string]interface{}{"email": request.Email}); err != nil || count != 0 {
		return utils.Error(http.StatusBadRequest, "该邮箱已被注册")
	}

	if err := utils.VerifyEmailCode(request.Email, request.Code); err != nil {
		return utils.Error(http.StatusBadRequest, err.Error())
	}

	var newUser model.User
	newUser.Username = request.Username
	newUser.Salt = utils.GenerateSalt(16)
	newUser.Password = utils.HashPassword(request.Password, newUser.Salt)
	newUser.CreatedAt = time.Now().UnixMilli()
	newUser.Avatar = request.Avatar
	newUser.Email = request.Email

	if err := us.userRepo.Create(&newUser); err != nil {
		return utils.Error(http.StatusInternalServerError, "创建用户时发生错误")
	}

	return utils.Success(http.StatusCreated)
}

// UpdateUser 更新用户信息
func (us *User) UpdateUser(id uint64, fileds interface{}, request *utils.UpdateUserRequest) *utils.Response {
	if request.Email != "" {
		if count, err := us.userRepo.Count(map[string]interface{}{"email": request.Email}); err != nil || count != 0 {
			return utils.Error(http.StatusBadRequest, "该邮箱已被注册")
		}

		if request.Code == "" {
			return utils.Error(http.StatusBadRequest, "验证码不能为空")
		}

		if err := utils.VerifyEmailCode(request.Email, request.Code); err != nil {
			return utils.Error(http.StatusBadRequest, err.Error())
		}
	}

	if request.Username != "" {
		if count, err := us.userRepo.Count(map[string]interface{}{"username": request.Username}); err != nil || count != 0 {
			return utils.Error(http.StatusBadRequest, "该用户名已被注册")
		}
	}

	if err := us.userRepo.Update(map[string]interface{}{"id": id}, fileds, request); err != nil {
		return utils.Error(http.StatusInternalServerError, "更新用户失败")
	}

	return utils.Success(http.StatusOK)
}

// DeleteUser 根据用户 ID 删除用户
func (us *User) DeleteUser(id uint64) *utils.Response {
	if err := us.userRepo.Delete(map[string]interface{}{"id": id}); err != nil {
		return utils.Error(http.StatusInternalServerError, "删除用户失败")
	}
	return utils.Success(http.StatusOK)
}

// UpdateUserPassword 更新用户密码
func (us *User) UpdateUserPassword(id uint64, request utils.UpdatePasswordRequest) *utils.Response {
	var result struct {
		Email    string
		Salt     string
		Password string
	}
	if err := us.userRepo.Search(map[string]interface{}{"id": id}, 1, &result); err != nil {
		return utils.Error(http.StatusUnauthorized, "用户不存在")
	}

	if utils.HashPassword(request.Password, result.Salt) != result.Password {
		return utils.Error(http.StatusUnauthorized, "密码错误")
	}

	if err := utils.VerifyEmailCode(result.Email, request.Code); err != nil {
		return utils.Error(http.StatusBadRequest, err.Error())
	}

	salt := utils.GenerateSalt(16)
	password := utils.HashPassword(request.NewPassword, salt)

	if err := us.userRepo.Update(map[string]interface{}{"id": id}, []string{"salt", "password"}, map[string]interface{}{"salt": salt, "password": password}); err != nil {
		return utils.Error(http.StatusInternalServerError, "内部错误")
	}
	return utils.Success(http.StatusOK)
}

func (us *User) SendEmailVerification(request utils.SendEmailVerificationRequest) *utils.Response {
	code := utils.GenerateCode(6)
	global.Rdb.Set(global.Ctx, request.Email, code, time.Minute*time.Duration(config.AppConfig.Email.Expiration))
	if err := utils.SendEmailVerification(request.Email, code); err != nil {
		return utils.Error(http.StatusInternalServerError, "发送验证码失败")
	}

	return utils.Success(http.StatusOK)
}
