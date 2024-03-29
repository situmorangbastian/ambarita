package models

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
)

// Article ...
type Article struct {
	ID         string    `json:"id" bson:"_id"`
	Title      string    `json:"title" bson:"title" validate:"required"`
	Slug       string    `json:"slug" bson:"slug"`
	Content    string    `json:"content" bson:"content" validate:"required"`
	CreateTime time.Time `json:"created_time" bson:"created_time"`
	UpdateTime time.Time `json:"updated_time" bson:"updated_time"`
}

func (a Article) Validate() error {
	cv := CustomValidator{
		Validator: validator.New(),
	}

	return cv.Validate(a)
}

// ArticleUsecase represent the article's usecases contract
type ArticleUsecase interface {
	Fetch(ctx context.Context, cursor string, num int) ([]Article, string, error)
	Get(ctx context.Context, ID string) (Article, error)
	Store(ctx context.Context, article Article) (Article, error)
	Update(ctx context.Context, article Article) (Article, error)
	Delete(ctx context.Context, ID string) error
}

// ArticleRepository represent the article's repository contract
type ArticleRepository interface {
	Fetch(ctx context.Context, cursor string, num int) ([]Article, string, error)
	Get(ctx context.Context, ID string) (Article, error)
	Store(ctx context.Context, article Article) error
	Update(ctx context.Context, article Article) error
	Delete(ctx context.Context, ID string) error
}
