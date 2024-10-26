package router

import (
	"fmt"
	"videohub/internal/config"
	"videohub/internal/controller"
	"videohub/internal/repository"
	"videohub/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func RouterInit(r *gin.Engine, config *config.Config) {

	s := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/videohub?charset=utf8mb4&parseTime=True&loc=Local", config.Mysql.Username, config.Mysql.Password)
	db, _ := gorm.Open(mysql.Open(s), &gorm.Config{})

	//依赖注入

	//1、db 到 repository
	collectionRepo := repository.NewCollection(db)
	commentRepo := repository.NewComment(db)
	userRepo := repository.NewUser(db)
	videoRepo := repository.NewVideo(db)

	//2、repository 到 service
	userAvatarService := service.NewUser_avatar(userRepo)
	userListService := service.NewUser_list(userRepo)
	userService := service.NewUser(userRepo, collectionRepo, videoRepo)
	videoListService := service.NewVideo_list(userRepo, videoRepo)
	videoUploadService := service.NewVideo_upload(userRepo, videoRepo)
	videoService := service.NewVideo(userRepo, videoRepo, commentRepo)

	//3、service 到 controller
	users_controller := controller.NewUsers(userAvatarService, userListService, userService)
	videos_controller := controller.NewVideos(videoService, videoListService, videoUploadService)

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
