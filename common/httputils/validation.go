package httputils

import "github.com/go-playground/validator/v10"

func Validate(T any) error {
	v := validator.New()
	return v.Struct(T)
}
