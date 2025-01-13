package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"videohub/config"
	"videohub/global"
	"videohub/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

// AuthMiddleware 认证中间件, role 为 0 表示普通用户, 1 表示管理员, 进行 token 验证
func AuthMiddleware(role int8) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析 token
		token := c.GetHeader("Authorization")
		if token == "" {
			logrus.Debug("token is invalid")
			c.JSON(http.StatusOK, utils.Error(http.StatusUnauthorized, "未授权"))
			c.Abort()
			return
		}
		payload, err := utils.ParseJWT(token, config.AppConfig.JWT.AccessTokenSecret)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				logrus.Debug(err.Error())
				c.JSON(http.StatusOK, utils.Error(http.StatusUnauthorized, "未授权"))
			} else {
				logrus.Error(err.Error())
				c.JSON(http.StatusOK, utils.Error(http.StatusUnauthorized, "未授权"))
			}
			c.Abort()
			return
		}

		// 如果是普通用户, 则设置用户在线状态
		if payload.Role == 0 {
			key := fmt.Sprintf("user:%d:is_online", payload.ID)
			global.Rdb.Set(global.Ctx, key, true, 1*time.Minute)
		}

		// 判断用户角色是否匹配
		if payload.Role == role {
			c.Set("id", payload.ID)
			c.Set("role", payload.Role)
			c.Next()
			return
		}
		c.JSON(http.StatusOK, utils.Error(http.StatusUnauthorized, "未授权"))
		c.Abort()
	}
}
