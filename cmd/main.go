package main

import (
	"fmt"
	"os"
	"videohub/internal/api/users"
	"videohub/internal/api/videos"

	"github.com/gin-gonic/gin"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Run     runConfig     `yaml:"run"`
	Storage storageConfig `yaml:"storage"`
}
type runConfig struct {
	IP   string `yaml:"ip"`
	Port string `yaml:"port"`
}

type storageConfig struct {
	Videos_data  string `yaml:"videos_data"`
	Videos_cover string `yaml:"videos_cover"`
	Images       string `yaml:"images"`
}

func main() {
	r := gin.Default()

	//路由初始化，其中包含工厂方法
	routerInit(r)

	dataBytes, err := os.ReadFile("config/application.yaml")
	if err != nil {
		fmt.Println("读取文件失败：", err)
		return
	}
	config := Config{}
	err = yaml.Unmarshal(dataBytes, &config)
	if err != nil {
		fmt.Println("解析 yaml 文件失败：", err)
		return
	}

	// 设置静态文件夹路径
	r.Static("/storage/images", config.Storage.Images)             // 图像存储
	r.Static("/storage/videos", config.Storage.Videos_data)        // 视频存储
	r.Static("/storage/videos_cover", config.Storage.Videos_cover) //视频封面存储

	r.Run(config.Run.IP + ":" + config.Run.Port)
}

func routerInit(r *gin.Engine) {

	// 用户路由组
	userRouter := r.Group("/users")
	{
		userRouter.GET("/", users.User_test)
		// 可以添加更多用户相关的路由
	}
	//视频路由组
	videoRouter := r.Group("/videos")
	{
		videoRouter.GET("/test", videos.Video_test)
	}
}
