package validator

import (
	validatorV10 "github.com/go-playground/validator/v10"
)

var validator = validatorV10.New()

func Validate(data any) error {
	return validator.Struct(data)
}
