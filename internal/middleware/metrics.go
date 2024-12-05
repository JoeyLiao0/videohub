package middleware

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"videohub/config"
	"videohub/global"
	"videohub/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CountViewMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.Next()
			return
		}
		payload, err := utils.ParseJWT(token, config.AppConfig.JWT.AccessTokenSecret)
		if err != nil {
			c.Next()
			return
		}
		id := payload.ID
		path := c.Request.URL.Path
		if c.Request.Method == http.MethodGet && strings.HasPrefix(path, "/static/videos/data/") {
			filename := filepath.Base(path)
			vid := strings.TrimSuffix(filename, filepath.Ext(filename))
			viewKey := "video:" + vid + ":views"
			if global.Rdb.Exists(global.Ctx, viewKey).Val() == 0 {
				c.Next()
				return
			}
			userViewKey := "user:" + strconv.Itoa(int(id)) + ":" + vid
			if global.Rdb.Exists(global.Ctx, userViewKey).Val() > 0 {
				c.Next()
				return
			}
			// ipViewKey := "ip:" + c.ClientIP() + ":" + vid
			// if global.Rdb.Exists(global.Ctx, ipViewKey).Val() > 0 {
			// 	c.Next()
			// 	return
			// }
			if err := global.Rdb.Incr(global.Ctx, viewKey).Err(); err != nil {
				logrus.Debug(err.Error())
				c.JSON(http.StatusOK, utils.Error(http.StatusInternalServerError, "服务器内部错误"))
				c.Abort()
			}
			global.Rdb.Set(global.Ctx, userViewKey, 1, 1*time.Hour)
			// global.Rdb.Set(global.Ctx, ipViewKey, 1, 5*time.Minute)
		}
		c.Next()
	}
}
