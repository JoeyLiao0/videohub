package config

import (
	"context"
	"log"
	"videohub/global"

	"github.com/redis/go-redis/v9"
)

func initRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     AppConfig.Redis.Host + ":" + AppConfig.Redis.Port,
		Password: AppConfig.Redis.Password,
		DB:       AppConfig.Redis.DB,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Error connecting to redis: %v", err)
	}
	log.Println("Connected to redis")

	global.Rdb = rdb
}
