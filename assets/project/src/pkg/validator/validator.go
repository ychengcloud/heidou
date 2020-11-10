package validator

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	Validate *validator.Validate
}

// Init 初始化数据库
func New() *Validator {
	validate := validator.New()
	return &Validator{
		Validate: validate,
	}
}
