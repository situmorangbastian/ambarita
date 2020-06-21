package models

import (
	"context"
	"time"
)

// Article ...
type Article struct {
	ID          string    `json:"id"`
	Title       string    `json:"title" validate:"required"`
	Slug        string    `json:"slug"`
	Content     string    `json:"content" validate:"required"`
	CreateTime  time.Time `json:"created_time"`
	UpdateTime  time.Time `json:"updated_time"`
	DeletedTime time.Time `json:"-"`
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
