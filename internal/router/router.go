package router

import (
	"videohub/config"
	"videohub/global"
	"videohub/internal/controller"
	"videohub/internal/middleware"
	"videohub/internal/repository"
	"videohub/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	db := global.DB

	// 1、db 到 repository
	collectionRepo := repository.NewCollection(db)
	userRepo := repository.NewUser(db)
	videoRepo := repository.NewVideo(db)
	commentRepo := repository.NewComment(db)
	likeRepo := repository.NewLike(db)
	statsRepo := repository.NewStats(db)

	// 2、repository 到 service
	userAvatarService := service.NewUserAvatar(userRepo)
	userListService := service.NewUserList(userRepo)
	userService := service.NewUser(userRepo, collectionRepo, videoRepo)
	videoUploadService := service.NewVideoUpload(videoRepo)
	VideoUpdateStatusService := service.NewVideoUpdateStatus(videoRepo)
	videoSearchService := service.NewVideoSearch(videoRepo, likeRepo, collectionRepo)
	commentService := service.NewComment(commentRepo, videoRepo)
	userVideoService := service.NewUserVideo(videoRepo, likeRepo, collectionRepo)
	userCollectionService := service.NewUserCollection(videoRepo, likeRepo, collectionRepo)
	likeService := service.NewLike(videoRepo, likeRepo)
	videoUpdateService := service.NewVideoUpdateStatus(videoRepo)
	statsService := service.NewStats(statsRepo)

	// 3、service 到 controller
	userController := controller.NewUserController(userAvatarService, userListService, userService, userVideoService, userCollectionService)
	videoController := controller.NewVideoController(videoUploadService, VideoUpdateStatusService, videoSearchService, likeService, commentService)
	adminController := controller.NewAdminController(
		userAvatarService,
		userListService,
		userService,
		videoSearchService,
		videoUpdateService,
		userVideoService,
		statsService,
	)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// 日志中间件
	r.Use(middleware.LoggerMiddleware())
	// CORS 跨域中间件
	r.Use(middleware.CORSMiddleware())

	// 静态文件路由组
	staticGroup := r.Group("/static")
	{
		// 设置静态文件夹路径  (url前缀, 文件夹路径)
		staticGroup.Static(config.AppConfig.Static.Avatar, config.AppConfig.Storage.Images)     // 头像存储
		staticGroup.Static(config.AppConfig.Static.Cover, config.AppConfig.Storage.VideosCover) // 视频封面存储
		staticGroup.Use(middleware.CountViewMiddleware())
		{
			staticGroup.Static(config.AppConfig.Static.Video, config.AppConfig.Storage.VideosData) // 视频存储
		}
	}

	// 管理员路由组
	adminRouter := r.Group("/admin")
	{
		// 管理员登录
		adminRouter.POST("/token", func(c *gin.Context) {
			userController.Login(c, 1)
		})
		// 利用刷新令牌获取访问令牌
		adminRouter.POST("/access_token", func(c *gin.Context) {
			userController.AccessToken(c, 1)
		})
		adminRouter.Use(middleware.AuthMiddleware(1))
		{
			// 获取管理员个人信息
			adminRouter.GET("", adminController.GetUser)
			// 获取用户信息
			adminRouter.GET("/users", adminController.GetUsers)
			// 创建用户
			adminRouter.POST("/users", adminController.CreateUser)
			// 更新用户信息
			adminRouter.PUT("/users", adminController.UpdateUser)

			// 视频列表获取
			adminRouter.GET("/videos", adminController.GetVideos)
			// 视频状态修改
			adminRouter.PUT("/videos", adminController.UpdateVideo)
			// 视频删除
			adminRouter.DELETE("/videos", adminController.DeleteVideo)
			// 管理员密码修改
			adminRouter.PUT("/password", userController.UpdatePassword)

			// 获取实时数据
			adminRouter.GET("/real_time_data", adminController.GetRealTimeData)
			// 获取历史数据
			adminRouter.GET("/historical_data", adminController.GetHistoricalData)
		}
	}

	// 用户路由组
	userRouter := r.Group("/users")
	{
		// 用户注册
		userRouter.POST("", userController.CreateUser)
		// 用户登录
		userRouter.POST("/token", func(c *gin.Context) {
			userController.Login(c, 0)
		})
		// 利用刷新令牌获取访问令牌
		userRouter.POST("/access_token", func(c *gin.Context) {
			userController.AccessToken(c, 0)
		})

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
			userRouter.DELETE("/videos", userController.DeleteVideo)
			// 获取用户视频收藏列表
			userRouter.GET("/collections", userController.GetCollections)
			// 用户收藏视频
			userRouter.POST("/collections", userController.AddCollection)
			// 用户删除收藏视频
			userRouter.DELETE("/collections", userController.DeleteCollection)
		}
	}

	// 视频路由组
	videoRouter := r.Group("/videos")
	{
		// 获取视频列表
		videoRouter.GET("", videoController.GetVideos)
		// 获取视频评论
		videoRouter.GET("/comments", videoController.GetComments)

		videoRouter.Use(middleware.AuthMiddleware(0))
		{
			// 更新视频状态
			videoRouter.PUT("", videoController.UpdateVideoStatus)
			// 视频点赞
			videoRouter.POST("/likes", videoController.LikeVideo)
			// 视频取消点赞
			videoRouter.DELETE("/likes", videoController.UnlikeVideo)
			// 新增视频评论
			videoRouter.POST("/comments", videoController.AddComment)
			// 删除评论
			videoRouter.DELETE("/comments", videoController.DeleteComment)
			// 评论点赞
			videoRouter.POST("/comments/likes", videoController.LikeComment)
			// 评论取消点赞
			videoRouter.DELETE("/comments/likes", videoController.UnlikeComment)
			// 视频分片上传
			videoRouter.POST("/chunk", videoController.UploadChunk)
			// 合并视频分片
			videoRouter.POST("/complete", videoController.CompleteUpload)
		}
	}

	// API 路由组
	apiRouter := r.Group("/api")
	{
		// 发送邮箱验证码
		apiRouter.POST("/email", userController.SendEmailVerification)
	}

	logrus.Info("router initialized successfully")
	return r
}
