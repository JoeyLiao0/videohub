package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{"https://localhost:8080"}, // 允许特定域
		// AllowMethods:     []string{"GET", "POST", "PUT", "OPTIONS"}, // 允许特定方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 允许特定请求头
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		//   return origin == "https://github.com"
		// },
		MaxAge: 12 * time.Hour,
	})
}