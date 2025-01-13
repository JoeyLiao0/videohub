package global

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// 定义全局变量
var (
	// Gorm 数据库连接
	DB       *gorm.DB
	// Redis 连接
	Rdb      *redis.Client
	// Validator 验证器
	Validate *validator.Validate
	// 全局上下文
	Ctx      context.Context
)
