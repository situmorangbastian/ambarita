package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/situmorangbastian/ambarita/models"
)

type handler struct {
	usecase models.ArticleUsecase
}

// NewHandler will initialize the articles/ resources endpoint
func NewHandler(e *echo.Echo, usecase models.ArticleUsecase) {
	handler := &handler{
		usecase: usecase,
	}

	e.GET("/articles", handler.fetch)
	e.POST("/articles", handler.store)
	e.PUT("/articles/:id", handler.update)
	e.GET("/articles/:id", handler.get)
	e.DELETE("/articles/:id", handler.delete)
}

func (h handler) fetch(c echo.Context) error {
	num := 0
	numStr := c.QueryParam("num")

	if numStr != "" {
		var err error
		num, err = strconv.Atoi(c.QueryParam("num"))
		if err != nil {
			return models.ErrBadRequest
		}
	}

	articles, nextCursor, err := h.usecase.Fetch(c.Request().Context(), c.QueryParam("cursor"), num)
	if err != nil {
		return err
	}

	response := models.DefaultSuccessResponse()
	response.Data = articles
	response.PageInfo = &map[string]interface{}{
		"next_cursor": nextCursor,
	}

	return c.JSON(http.StatusOK, response)
}

func (h handler) get(c echo.Context) error {
	article, err := h.usecase.Get(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}

	response := models.DefaultSuccessResponse()
	response.Data = article

	return c.JSON(http.StatusOK, response)
}

func (h handler) store(c echo.Context) error {
	article := models.Article{}
	err := c.Bind(&article)
	if err != nil {
		return models.ErrBadRequest
	}

	if err := c.Validate(article); err != nil {
		return err
	}

	storedArticle, err := h.usecase.Store(c.Request().Context(), article)
	if err != nil {
		return err
	}

	response := models.DefaultCreatedResponse()
	response.Data = storedArticle

	return c.JSON(http.StatusCreated, response)
}

func (h handler) update(c echo.Context) error {
	article := models.Article{}
	err := c.Bind(&article)
	if err != nil {
		return models.ErrBadRequest
	}

	article.ID = c.Param("id")

	if err := c.Validate(article); err != nil {
		return err
	}

	updatedArticle, err := h.usecase.Update(c.Request().Context(), article)
	if err != nil {
		return err
	}

	response := models.DefaultSuccessResponse()
	response.Data = updatedArticle

	return c.JSON(http.StatusCreated, updatedArticle)
}

func (h handler) delete(c echo.Context) error {
	err := h.usecase.Delete(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
