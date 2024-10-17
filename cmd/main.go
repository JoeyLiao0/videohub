/*
*@auther:廖嘉鹏
*项目的启动文件
 */

package main

import (
	"fmt"
	"videohub/internal/api"
	"videohub/internal/repository"
	"videohub/internal/service"
	"videohub/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	config, err := utils.GetConfig()

	if err != nil || config == nil {
		fmt.Println("读取配置文件/解析配置文件失败：", err)
		return
	}

	r := gin.Default()

	//路由初始化，其中包含工厂方法
	routerInit(r, config)

	// 设置静态文件夹路径
	r.Static("/storage/images", config.Storage.Images)             // 图像存储
	r.Static("/storage/videos", config.Storage.Videos_data)        // 视频存储
	r.Static("/storage/videos_cover", config.Storage.Videos_cover) //视频封面存储

	r.Run(config.Run.IP + ":" + config.Run.Port)
}

func routerInit(r *gin.Engine, config *utils.Config) {

	s := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/videohub?charset=utf8mb4&parseTime=True&loc=Local", config.Mysql.Username, config.Mysql.Password)
	db, _ := gorm.Open(mysql.Open(s), &gorm.Config{})

	//依赖注入

	//1、db 到 repository
	collection_repository := repository.NewCollection_respository(db)
	comment_repository := repository.NewComment_respository(db)
	user_respository := repository.NewUser_respository(db)
	video_repository := repository.NewVideo_respository(db)

	//2、repository 到 service
	user_avatar_service := service.NewUser_avatar_service(user_respository)
	user_list_service := service.NewUser_list_service(user_respository)
	user_service := service.NewUser_service(user_respository, collection_repository, video_repository)
	video_list_service := service.NewVideo_list_service(user_respository, video_repository)
	video_upload_service := service.NewVideo_upload_service(user_respository, video_repository)
	video_service := service.NewVideo_service(user_respository, video_repository, comment_repository)

	//3、service 到 controller
	users_controller := api.NewUsers_controller(user_avatar_service, user_list_service, user_service)
	videos_controller := api.NewVideos_controller(video_service, video_list_service, video_upload_service)

	//完全注入好了,进行路由

	// 用户路由组
	userRouter := r.Group("/users")
	{
		userRouter.GET("/test", users_controller.User_test)
		// 可以添加更多用户相关的路由,一个API对应一个控制层函数
	}
	//视频路由组
	videoRouter := r.Group("/videos")
	{
		videoRouter.GET("/test", videos_controller.Video_test)
		// 可以添加更多用户相关的路由,一个API对应一个控制层函数
	}
}
