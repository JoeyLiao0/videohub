package config

import (
	"videohub/global"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func initRedis() {
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
	global.Rdb = rdb
}
