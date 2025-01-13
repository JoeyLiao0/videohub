package config

import (
	"videohub/global"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// InitRedis 初始化 Redis
func InitRedis() {
	// 连接 Redis(client)
	rdb := redis.NewClient(&redis.Options{
		Addr:     AppConfig.Redis.Host + ":" + AppConfig.Redis.Port,
		Password: AppConfig.Redis.Password,
		DB:       AppConfig.Redis.DB,
	})
	_, err := rdb.Ping(global.Ctx).Result()
	if err != nil {
		logrus.Fatalf("Error connecting to redis: %v", err)
	}
	logrus.Info("Redis connected successfully")

	// 可以加入其他的 Redis 配置
	global.Rdb = rdb
}
