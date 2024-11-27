package middleware

import (
	"time"
	"videohub/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     config.AppConfig.CORS.AllowOrigins,     // 允许所有域名
		AllowMethods:     config.AppConfig.CORS.AllowMethods,     // 允许特定方法
		AllowHeaders:     config.AppConfig.CORS.AllowHeaders,     // 允许特定请求头
		ExposeHeaders:    config.AppConfig.CORS.ExposeHeaders,    // 允许特定响应头
		AllowCredentials: config.AppConfig.CORS.AllowCredentials, // 允许携带凭证
		// AllowOriginFunc: func(origin string) bool {
		//   return origin == "http://localhost:8080"
		// },
		MaxAge: time.Duration(config.AppConfig.CORS.MaxAge) * time.Second, // 预检请求有效期
	})
}
