package router

import (
	"videohub/config"
	"videohub/internal/controller"
	"videohub/internal/middleware"
	"videohub/internal/repository"
	"videohub/internal/service"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	db := config.InitDB()

	//1、db 到 repository
	collectionRepo := repository.NewCollection(db)
	userRepo := repository.NewUser(db)
	videoRepo := repository.NewVideo(db)

	//2、repository 到 service
	userAvatarService := service.NewUser_avatar(userRepo)
	userListService := service.NewUser_list(userRepo)
	userService := service.NewUser(userRepo, collectionRepo, videoRepo)
	videoUploadService := service.NewVideoUploadService(videoRepo)

	//3、service 到 controller
	users_controller := controller.NewUsers(userAvatarService, userListService, userService)
	videos_controller := controller.NewVideos(videoUploadService)

	r := gin.Default()
	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	// r.MaxMultipartMemory = 8 << 20  // 8 MiB

	// CORS 跨域中间件
	r.Use(middleware.CORSMiddleware())

	// 设置静态文件夹路径
	r.Static("/storage/images", config.AppConfig.Storage.Images)            // 图像存储
	r.Static("/storage/videos", config.AppConfig.Storage.VideosData)        // 视频存储
	r.Static("/storage/videos_cover", config.AppConfig.Storage.VideosCover) // 视频封面存储

	// 用户路由组
	userRouter := r.Group("/users")
	{
		userRouter.POST("", users_controller.CreateUser)

		userRouter.POST("/token", users_controller.Login)
		userRouter.POST("/access_token", users_controller.AccessToken)

		userRouter.Use(middleware.AuthMiddleware())
		{
			userRouter.GET("", users_controller.GetAllUsers)

			userRouter.GET("/:id", users_controller.GetUserByID)
			userRouter.PUT("/:id", users_controller.UpdateUser)
			userRouter.DELETE("/:id", users_controller.DeleteUser)

			userRouter.POST("/:id/avatar", users_controller.UploadAvatar)

			userRouter.PUT("/:id/password", users_controller.UpdatePassword)

			userRouter.POST("/:id/email", users_controller.SendEmailVerification)

		}
	}
	// 视频路由组
	videoRouter := r.Group("/videos")
	{
		// videoRouter.GET("", videos_controller.GetVideos)
		videoRouter.POST("/chunk", videos_controller.UploadChunk)
		videoRouter.POST("/complete", videos_controller.CompleteUpload)
	}

	return r
}
