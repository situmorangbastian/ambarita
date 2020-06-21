package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"github.com/situmorangbastian/ambarita/models"
)

// ContextTimeout returns custom middleware for Echo that set maximum HTTP response time
// before considered timeout with duration d.
func ContextTimeout(d time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, cancel := context.WithTimeout(c.Request().Context(), d)
			defer cancel()

			newRequest := c.Request().WithContext(ctx)
			c.SetRequest(newRequest)

			return next(c)
		}
	}
}

// ErrMiddleware returns custom middleware for Echo that generate HTTP error response
// with HTTP status code.
func ErrMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			// Handle Error from Echo
			errorEcho, ok := err.(*echo.HTTPError)
			if ok && errorEcho.Code < 500 {
				return nil
			}

			errResponse := models.DefaultErrorResponse()
			errResponse.Message = err.Error()

			// Check error based on error type
			switch err.(type) {
			case models.ErrorValidation:
				errResponse.Status = http.StatusBadRequest
				errResponse.Data = map[string]interface{}{}
				return echo.NewHTTPError(errResponse.Status, errResponse)
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

			return echo.NewHTTPError(errResponse.Status, errResponse)
		}
	}
}
