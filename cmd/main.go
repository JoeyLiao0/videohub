package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"
	"videohub/config"
	"videohub/global"
	"videohub/internal/repository"
	"videohub/internal/router"
	"videohub/internal/utils"
	"videohub/logger"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func main() {
	global.Ctx = context.Background()
	config.InitConfig()
	logger.InitLogger(config.AppConfig.Run.Debug)
	config.InitDB()
	config.InitRedis()
	utils.InitValidator()
	r := router.InitRouter()
	srv := &http.Server{
		Addr:    config.AppConfig.Run.Host + ":" + config.AppConfig.Run.Port,
		Handler: r,
	}

	c := cron.New()
	c.AddFunc("0 0 0 * * *", func() { repository.WriteStats(global.DB) })
	c.Start()
	// repository.WriteStats(global.DB)
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server Shutdown:", err)
	}
	logrus.Info("Server exiting")
	logrus.Info()
}
