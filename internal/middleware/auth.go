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

func AuthMiddleware(role int8) gin.HandlerFunc {
	return func(c *gin.Context) {
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
				c.JSON(http.StatusOK, utils.Error(http.StatusInternalServerError, "未授权"))
			}
			c.Abort()
			return
		}

		key := fmt.Sprintf("user:%d:is_online", payload.ID)
		if global.Rdb.Exists(global.Ctx, key).Val() > 0 {
			global.Rdb.Expire(global.Ctx, key, 1*time.Minute)
		} else {
			global.Rdb.Set(global.Ctx, key, true, 1*time.Minute)
		}

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
