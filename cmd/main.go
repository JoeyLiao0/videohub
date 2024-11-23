/*
*@auther:廖嘉鹏
*项目的启动文件
 */

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"videohub/config"
	"videohub/global"
	"videohub/internal/router"
	"videohub/internal/utils"
)

func main() {
	global.Ctx = context.Background()
	config.InitConfig()
	utils.InitValidator()
	r := router.InitRouter()
	// r.Run(config.AppConfig.Run.IP + ":" + config.AppConfig.Run.Port)
	srv := &http.Server{
		Addr:    config.AppConfig.Run.Host + ":" + config.AppConfig.Run.Port,
		Handler: r,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
