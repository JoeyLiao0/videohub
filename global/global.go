package global

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	DB       *gorm.DB
	Rdb      *redis.Client
	Validate *validator.Validate
	Ctx      context.Context
)
