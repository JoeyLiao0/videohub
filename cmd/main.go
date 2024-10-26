/*
*@auther:廖嘉鹏
*项目的启动文件
 */

package main

import (
	"videohub/internal/config"
	"videohub/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()
	config := config.GetConfig()

	if config == nil {
		return
	}

	r := gin.Default()

	//路由初始化，其中包含工厂方法
	router.RouterInit(r, config)

	// 设置静态文件夹路径
	r.Static("/storage/images", config.Storage.Images)             // 图像存储
	r.Static("/storage/videos", config.Storage.Videos_data)        // 视频存储
	r.Static("/storage/videos_cover", config.Storage.Videos_cover) //视频封面存储

	r.Run(config.Run.IP + ":" + config.Run.Port)
}
