package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"

	"github.com/situmorangbastian/ambarita/models"
)

// ErrMiddleware returns custom middleware for Fiber that generate HTTP error response
// with HTTP status code.
func ErrMiddleware(ctx *fiber.Ctx, err error) error {
	errResponse := models.DefaultErrorResponse()
	errResponse.Message = err.Error()

	// Retreive the custom statuscode if it's an fiber.*Error
	if e, ok := err.(*fiber.Error); ok {
		errResponse.Message = e.Error()
		errResponse.Status = e.Code
		return ctx.Status(errResponse.Status).JSON(errResponse)
	}

	// Check error based on error type
	switch err.(type) {
	case models.ErrorValidation:
		errResponse.Status = http.StatusBadRequest
		errResponse.Data = map[string]interface{}{}
		return ctx.Status(errResponse.Status).JSON(errResponse)
	}

	switch err {
	case models.ErrBadRequest:
		errResponse.Status = http.StatusBadRequest
	case models.ErrNotFound:
		errResponse.Status = http.StatusNotFound
	default:
		log.Error(err)
		errResponse.Message = "Internal Server Error"
	}

	return ctx.Status(errResponse.Status).JSON(errResponse)
}
