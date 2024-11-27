package utils

import (
	"videohub/global"

	"github.com/go-playground/validator/v10"
)

func InitValidator() {
	global.Validate = validator.New(validator.WithRequiredStructEnabled())
}