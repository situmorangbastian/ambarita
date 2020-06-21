package models

import (
	"strings"

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
		errorValidators, ok := err.(validator.ValidationErrors)
		if ok {
			errMessages := []string{}
			for _, errValidator := range errorValidators {
				errMessages = append(errMessages, errValidator.Field()+" "+errValidator.ActualTag())
			}
			err = NewErrorValidation(strings.Join(errMessages, ", "))
		}
	}

	return err
}
