package http

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	gowerModel "github.com/situmorangbastian/gower/models"

	"github.com/situmorangbastian/ambarita/models"
)

type handler struct {
	usecase models.ArticleUsecase
}

// NewHandler will initialize the articles/ resources endpoint
func NewHandler(f *fiber.App, usecase models.ArticleUsecase) {
	handler := &handler{
		usecase: usecase,
	}

	f.Get("/articles", handler.fetch)
	f.Post("/articles", handler.store)
	f.Put("/articles/:id", handler.update)
	f.Get("/articles/:id", handler.get)
	f.Delete("/articles/:id", handler.delete)
}

func (h handler) fetch(c *fiber.Ctx) error {
	num := 0
	numStr := c.Query("num")

	if numStr != "" {
		var err error
		num, err = strconv.Atoi(c.Query("num"))
		if err != nil {
			return gowerModel.ConstraintErrorf("invalid query param num")
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
		return gowerModel.ConstraintErrorf(err.Error())
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
		return gowerModel.ConstraintErrorf(err.Error())
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

	return c.SendStatus(http.StatusNoContent)
}
