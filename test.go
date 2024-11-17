// package main

// import (
// 	"log"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	r := gin.Default()
// 	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
// 		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
// 	}

// 	r.POST("/foo", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, "foo")
// 	})

// 	r.GET("/bar", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, "bar")
// 	})

// 	r.GET("/status", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, "ok")
// 	})

// 	// 监听并在 0.0.0.0:8080 上启动服务
// 	r.Run()
// }

package main

import "fmt"

func main() {
    var role uint8 = 0
    fmt.Println(role == 0) // 应该输出 true
}