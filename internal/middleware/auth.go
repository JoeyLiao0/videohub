package middleware

import (
	"errors"
	"log"
	"net/http"
	"videohub/config"
	"videohub/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		log.Println("token:", token)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})  // 401
			c.Abort()
			return
		}
		payload, err := utils.ParseJWT(token, config.AppConfig.JWT.AccessTokenSecret)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()}) // 401
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"}) // 500
			}
			c.Abort()
			return
		}
		c.Set("id", payload.ID)
		c.Set("role", payload.Role)
		c.Next()
	}
}