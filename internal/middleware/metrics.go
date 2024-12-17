package middleware

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"videohub/global"
	"videohub/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CountViewMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if c.Request.Method == http.MethodGet && strings.HasPrefix(path, "/static/videos/data/") {
			filename := filepath.Base(path)
			vid := strings.TrimSuffix(filename, filepath.Ext(filename))
			viewKey := "video:" + vid + ":views"
			if global.Rdb.Exists(global.Ctx, viewKey).Val() == 0 {
				c.Next()
				return
			}
			ipViewKey := "ip:" + c.ClientIP() + ":" + vid
			if global.Rdb.Exists(global.Ctx, ipViewKey).Val() > 0 { 
				c.Next()
				return
			}
			if err := global.Rdb.Incr(global.Ctx, viewKey).Err(); err != nil {
				logrus.Debug(err.Error())
				c.JSON(http.StatusOK, utils.Error(http.StatusInternalServerError, "服务器内部错误"))
				c.Abort()
			}
			global.Rdb.Set(global.Ctx, ipViewKey, 1, 1*time.Hour)
		}
		c.Next()
	}
}
