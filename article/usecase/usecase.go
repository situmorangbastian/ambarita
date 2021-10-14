package usecase

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/situmorangbastian/ambarita/models"
	"github.com/situmorangbastian/eclipse"
)

type usecase struct {
	repository models.ArticleRepository
}

// NewArticleUsecase will create new an usecase object representation of models.ArticleUsecase interface
func NewArticleUsecase(repository models.ArticleRepository) models.ArticleUsecase {
	return &usecase{
		repository: repository,
	}
}

func (u usecase) Fetch(ctx context.Context, cursor string, num int) ([]models.Article, string, error) {
	articles, nextCursor, err := u.repository.Fetch(ctx, cursor, num)
	if err != nil {
		return make([]models.Article, 0), "", errors.Wrap(err, "fetch article")
	}

	if len(articles) == 0 {
		nextCursor = cursor
	}

	return articles, nextCursor, nil
}

func (u usecase) Get(ctx context.Context, ID string) (models.Article, error) {
	article, err := u.repository.Get(ctx, ID)
	if err != nil {
		return models.Article{}, errors.Wrap(err, "get article")
	}

	return article, nil
}

func (u usecase) Store(ctx context.Context, article models.Article) (models.Article, error) {
	article.ID = uuid.New().String()
	article.CreateTime = time.Now()
	article.UpdateTime = time.Now()

	slug, err := u.resolveSlug(ctx, buildSlug(article.Title))
	if err != nil {
		return models.Article{}, errors.Wrap(err, "resolve slug")
	}
	article.Slug = slug

	err = u.repository.Store(ctx, article)
	if err != nil {
		return models.Article{}, err
	}

	return article, nil
}

func (u usecase) Update(ctx context.Context, article models.Article) (models.Article, error) {
	currentArticle, err := u.Get(ctx, article.ID)
	if err != nil {
		return models.Article{}, errors.Wrap(err, "get article on update")
	}

	article.Slug = currentArticle.Slug
	article.UpdateTime = time.Now()

	err = u.repository.Update(ctx, article)
	if err != nil {
		return models.Article{}, errors.Wrap(err, "update article")
	}

	article.CreateTime = currentArticle.CreateTime

	return article, nil
}

func (u usecase) Delete(ctx context.Context, id string) error {
	err := u.repository.Delete(ctx, id)
	if err != nil {
		return errors.Wrap(err, "delete article")
	}
	return nil
}

func (u usecase) resolveSlug(ctx context.Context, slug string) (string, error) {
	_, err := u.repository.Get(ctx, slug)
	if err != nil {
		switch errors.Cause(err).(type) {
		case eclipse.NotFoundError:
			return slug, nil
		default:
			return "", err
		}
	}

	counterSlug := int(1)
	for {
		newSlug := slug + "-" + strconv.Itoa(counterSlug)
		_, err = u.repository.Get(ctx, newSlug)
		if err != nil {
			switch errors.Cause(err).(type) {
			case eclipse.NotFoundError:
				return newSlug, nil
			default:
				return "", err
			}
		}

		counterSlug++
	}
}

func buildSlug(keyword string) string {
	regex := regexp.MustCompile("[^a-zA-Z0-9]+")
	processedSlug := regex.ReplaceAllString(keyword, " ")
	toLowerCase := strings.ToLower(processedSlug)
	splitTitle := strings.Fields(toLowerCase)
	slug := strings.Join(splitTitle, "-")
	return slug
}
