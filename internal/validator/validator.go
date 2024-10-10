package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Init() {
	validate = validator.New()
}

func Get() *validator.Validate {
	return validate
}
