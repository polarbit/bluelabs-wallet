package api

import (
	"github.com/go-playground/validator"
)

type CustomEchoValidator struct {
	v *validator.Validate
}

func (v *CustomEchoValidator) Init() {

}

func (v *CustomEchoValidator) Validate(i interface{}) error {
	return v.v.Struct(i)
}
