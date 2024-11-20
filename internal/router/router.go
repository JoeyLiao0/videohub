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
	userAvatarService := service.NewUserAvatar(userRepo)
	userListService := service.NewUserList(userRepo)
	userService := service.NewUser(userRepo, collectionRepo, videoRepo)
	videoUploadService := service.NewVideoUploadService(videoRepo)

	//3、service 到 controller
	userController := controller.NewUserController(userAvatarService, userListService, userService)
	videoController := controller.NewVideoController(videoUploadService)
	adminController := controller.NewAdminController(userAvatarService, userListService, userService)

	r := gin.Default()
	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	// r.MaxMultipartMemory = 8 << 20  // 8 MiB

	// CORS 跨域中间件
	r.Use(middleware.CORSMiddleware())

	// 设置静态文件夹路径
	r.Static("/storage/images", config.AppConfig.Storage.Images)            // 图像存储
	r.Static("/storage/videos", config.AppConfig.Storage.VideosData)        // 视频存储
	r.Static("/storage/videos_cover", config.AppConfig.Storage.VideosCover) // 视频封面存储

	adminRouter := r.Group("/admin")
	{
		adminRouter.POST("/token", userController.Login)
		adminRouter.POST("/access_token", userController.AccessToken)
		adminRouter.Use(middleware.AuthMiddleware(1))
		{
			// 获取管理员个人信息
			adminRouter.GET("", adminController.GetUser)
			// 获取用户信息
			adminRouter.GET("/users", adminController.GetUsers)
			// 创建用户
			adminRouter.POST("/users", adminController.CreateUser)
			// 更新用户信息
			adminRouter.PUT("/users/:id", adminController.UpdateUser)

			// 视频列表获取
			adminRouter.GET("/videos", adminController.GetVideos)
			// 视频状态修改
			adminRouter.PUT("/videos/:vid", adminController.UpdateVideo)
			// 视频删除
			adminRouter.DELETE("/videos/:vid", adminController.DeleteVideo)
		}
	}

	// 用户路由组
	userRouter := r.Group("/users")
	{
		// 用户注册
		userRouter.POST("", userController.CreateUser)
		// 用户登录
		userRouter.POST("/token", userController.Login)
		// 利用刷新令牌获取访问令牌
		userRouter.POST("/access_token", userController.AccessToken)

		userRouter.Use(middleware.AuthMiddleware(0))
		{
			// 获取用户信息
			userRouter.GET("", userController.GetUser)
			// 修改用户信息
			userRouter.PUT("", userController.UpdateUser)
			// 删除用户 (注销, 是否从数据库删除, 还是只修改 status)
			userRouter.DELETE("", userController.DeleteUser)
			// 用户上传头像
			userRouter.POST("/avatar", userController.UploadAvatar)
			// 用户修改密码
			userRouter.PUT("/password", userController.UpdatePassword)
			// 用户发布视频列表获取
			userRouter.GET("/videos", userController.GetVideos)
			// 删除用户发布的视频
			userRouter.DELETE("/videos/:vid", userController.DeleteVideo)
			// 获取用户视频收藏列表
			userRouter.GET("/collections", userController.GetCollections)
			// 用户收藏视频
			userRouter.POST("/collections", userController.UpdateCollections)
			// 用户删除收藏视频
			userRouter.DELETE("/collections", userController.DeleteCollections)
		}
	}
	// 视频路由组
	videoRouter := r.Group("/videos")
	{
		// 获取视频列表
		videoRouter.GET("", videoController.GetVideos)
		// 获取视频评论
		videoRouter.GET("/:vid/comments", videoController.GetComments)

		videoRouter.Use(middleware.AuthMiddleware(0))
		{
			// 视频点赞
			videoRouter.POST("/:vid", videoController.LikeVideo)
			// 新增视频评论
			videoRouter.POST("/:vid/comments", videoController.AddComment)
			// 评论点赞
			videoRouter.POST("/:vid/comments/:cid", videoController.LikeComment)
			// 删除评论
			videoRouter.DELETE("/:vid/comments/:cid", videoController.DeleteComment)
			// 视频分片上传
			videoRouter.POST("/chunk", videoController.UploadChunk)
			// 合并视频分片
			videoRouter.POST("/complete", videoController.CompleteUpload)
		}
	}

	apiRouter := r.Group("/api")
	{
		// 发送邮箱验证码
		apiRouter.POST("/email", userController.SendEmailVerification)
	}
	return r
}
