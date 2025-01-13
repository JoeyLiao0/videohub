package utils

import (
	"videohub/global"

	"github.com/go-playground/validator/v10"
)

// InitValidator 初始化验证器
func InitValidator() {
	global.Validate = validator.New(validator.WithRequiredStructEnabled())
}