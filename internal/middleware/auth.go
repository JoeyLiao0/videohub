package middleware

import (
	"errors"
	"net/http"
	"videohub/config"
	"videohub/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func AuthMiddleware(role int8) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token不存在"}) // 401
			c.Abort()
			return
		}
		payload, err := utils.ParseJWT(token, config.AppConfig.JWT.AccessTokenSecret)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token失效"}) // 401
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "内部错误"}) // 500
			}
			c.Abort()
			return
		}
		logrus.Debugf("role: %v %T", payload.Role, payload.Role)
		if payload.Role == role {
			c.Set("id", payload.ID)
			c.Set("role", payload.Role)
			c.Next()
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "禁止访问"}) // 403
		c.Abort()
	}
}
