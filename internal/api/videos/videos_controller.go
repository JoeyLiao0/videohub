package videos

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Video_test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "用户首页"})
}
