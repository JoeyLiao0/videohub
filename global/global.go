package global

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

var (
	Rdb      *redis.Client
	Validate *validator.Validate
	Ctx      context.Context
)
