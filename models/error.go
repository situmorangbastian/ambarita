package models

import "errors"

var (
	// ErrBadRequest will throw if the given request-body or params is not valid
	ErrBadRequest = errors.New("Given Param is not valid")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("Your requested Item is not found")
)

// ErrorValidation is a type specifically for error validation.
type ErrorValidation struct {
	message string
}

// Error retuns error message string.
func (err ErrorValidation) Error() string {
	return err.message
}

// NewErrorValidation is used to initialize error validation.
func NewErrorValidation(message string) ErrorValidation {
	return ErrorValidation{message: message}
}
