package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"

	"github.com/situmorangbastian/ambarita/models"
)

type handler struct {
	usecase models.ArticleUsecase
}

// NewHandler will initialize the articles/ resources endpoint
func NewHandler(f *fiber.App, usecase models.ArticleUsecase, requestTimeout time.Duration) {
	handler := &handler{
		usecase: usecase,
	}

	f.Get("/articles", timeout.New(handler.fetch, requestTimeout*time.Second))
	f.Post("/articles", timeout.New(handler.store, requestTimeout*time.Second))
	f.Put("/articles/:id", timeout.New(handler.update, requestTimeout*time.Second))
	f.Get("/articles/:id", timeout.New(handler.get, requestTimeout*time.Second))
	f.Delete("/articles/:id", timeout.New(handler.delete, requestTimeout*time.Second))
}

func (h handler) fetch(c *fiber.Ctx) error {
	num := 0
	numStr := c.Query("num")

	if numStr != "" {
		var err error
		num, err = strconv.Atoi(c.Query("num"))
		if err != nil {
			return models.ErrBadRequest
		}
	}

	articles, nextCursor, err := h.usecase.Fetch(c.Context(), c.Query("cursor"), num)
	if err != nil {
		return err
	}

	response := models.DefaultSuccessResponse()
	response.Data = articles
	response.PageInfo = &map[string]interface{}{
		"next_cursor": nextCursor,
	}

	return c.Status(http.StatusOK).JSON(response)
}

func (h handler) get(c *fiber.Ctx) error {
	article, err := h.usecase.Get(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}

	response := models.DefaultSuccessResponse()
	response.Data = article

	return c.Status(http.StatusOK).JSON(response)
}

func (h handler) store(c *fiber.Ctx) error {
	article := models.Article{}
	err := c.BodyParser(&article)
	if err != nil {
		return models.ErrBadRequest
	}

	if err := article.Validate(); err != nil {
		return err
	}

	storedArticle, err := h.usecase.Store(c.Context(), article)
	if err != nil {
		return err
	}

	response := models.DefaultCreatedResponse()
	response.Data = storedArticle

	return c.Status(http.StatusCreated).JSON(response)
}

func (h handler) update(c *fiber.Ctx) error {
	article := models.Article{}
	err := c.BodyParser(&article)
	if err != nil {
		return models.ErrBadRequest
	}

	article.ID = c.Params("id")

	if err := article.Validate(); err != nil {
		return err
	}

	updatedArticle, err := h.usecase.Update(c.Context(), article)
	if err != nil {
		return err
	}

	response := models.DefaultSuccessResponse()
	response.Data = updatedArticle

	return c.Status(http.StatusOK).JSON(updatedArticle)
}

func (h handler) delete(c *fiber.Ctx) error {
	err := h.usecase.Delete(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}

	return c.Status(http.StatusNoContent).Next()
}
