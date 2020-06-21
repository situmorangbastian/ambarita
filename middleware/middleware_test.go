package middleware

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	model "github.com/situmorangbastian/ambarita/models"
)

func Test_ContextTimeout(t *testing.T) {
	t.Run("error context timeout", func(t *testing.T) {
		mw := ContextTimeout(10 * time.Millisecond)

		h := func(c echo.Context) error {
			request := func(ctx context.Context) error {
				select {
				case <-time.After(100 * time.Millisecond):
					return nil
				case <-ctx.Done():
					err := ctx.Err()
					return err
				}
			}

			err := request(c.Request().Context())
			return err
		}

		req := httptest.NewRequest(http.MethodGet, "/publishers", nil)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		err := mw(h)(c)
		require.Error(t, err)
		require.Equal(t, context.DeadlineExceeded, err)
	})

	t.Run("success", func(t *testing.T) {
		mw := ContextTimeout(10 * time.Millisecond)

		h := func(c echo.Context) error {
			return nil
		}

		req := httptest.NewRequest(http.MethodGet, "/articles", nil)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		err := mw(h)(c)
		require.NoError(t, err)
	})
}

func Test_ErrMiddleware(t *testing.T) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)

	mw := ErrMiddleware()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	t.Run("error not found", func(t *testing.T) {
		h := func(c echo.Context) error {
			return model.ErrNotFound
		}

		err := mw(h)(c).(*echo.HTTPError)

		require.Error(t, err)
		require.Equal(t, http.StatusNotFound, err.Code)
		require.Contains(t, err.Error(), model.ErrNotFound.Error())
	})

	t.Run("error bad request", func(t *testing.T) {
		h := func(c echo.Context) error {
			return model.ErrBadRequest
		}

		err := mw(h)(c).(*echo.HTTPError)

		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.Code)
		require.Contains(t, err.Error(), model.ErrBadRequest.Error())
	})

	t.Run("error validation", func(t *testing.T) {
		h := func(c echo.Context) error {
			return model.NewErrorValidation("title is required")
		}

		err := mw(h)(c).(*echo.HTTPError)

		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.Code)
		require.Contains(t, err.Error(), model.NewErrorValidation("title is required").Error())
	})

	t.Run("internal server error", func(t *testing.T) {
		h := func(c echo.Context) error {
			return errors.New("unexpected error")
		}

		buf := new(bytes.Buffer)
		log.SetOutput(buf)

		err := mw(h)(c).(*echo.HTTPError)

		require.Error(t, err)
		require.Equal(t, http.StatusInternalServerError, err.Code)
		require.Contains(t, err.Error(), errors.New("Internal Server Error").Error())
		require.Contains(t, buf.String(), errors.New("unexpected error").Error())
	})
}
