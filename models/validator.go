package models

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator is custom validation for request data
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate validates given i based on validate struct tag
func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.Validator.Struct(i)
	if err != nil {
		for _, errValidator := range err.(validator.ValidationErrors) {
			return NewErrorValidation(errValidator.Field() + " " + errValidator.ActualTag())
		}
	}

	return err
}
