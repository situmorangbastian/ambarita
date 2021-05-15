package repository

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	log "github.com/sirupsen/logrus"
	"github.com/situmorangbastian/gower"

	"github.com/situmorangbastian/ambarita/models"
)

type mysqlRepository struct {
	db *sql.DB
}

// NewMysqlRepository will create the article mysql repository
func NewMysqlRepository(db *sql.DB) models.ArticleRepository {
	return mysqlRepository{
		db: db,
	}
}

func (r mysqlRepository) Fetch(ctx context.Context, cursor string, num int) ([]models.Article, string, error) {
	qBuilder := sq.Select("id", "title", "slug", "content", "created_time", "updated_time").
		From("articles").
		Where("deleted_time IS NULL").
		OrderBy("created_time DESC")

	if num > 0 {
		qBuilder = qBuilder.Limit(uint64(num))
	}

	if cursor != "" {
		decodedCursor, err := DecodeCursor(cursor)
		if err != nil {
			return make([]models.Article, 0), "", gower.ConstraintErrorf("invalid query param cursor")
		}
		qBuilder = qBuilder.Where(sq.Lt{"created_time": decodedCursor})
	}

	query, args, err := qBuilder.ToSql()
	if err != nil {
		return []models.Article{}, "", err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return []models.Article{}, "", err
	}

	articles := make([]models.Article, 0)
	for rows.Next() {
		article := models.Article{}
		err = rows.Scan(
			&article.ID,
			&article.Title,
			&article.Slug,
			&article.Content,
			&article.CreateTime,
			&article.UpdateTime,
		)
		if err != nil {
			log.Error(err)
			continue
		}

		articles = append(articles, article)
	}

	nextCursor := ""
	if len(articles) > 0 {
		nextCursor = EncodeCursor(articles[len(articles)-1].CreateTime)
	}

	return articles, nextCursor, nil
}

func (r mysqlRepository) Get(ctx context.Context, ID string) (models.Article, error) {
	query, args, err := sq.Select("id", "title", "slug", "content", "created_time", "updated_time").
		From("articles").
		Where(sq.And{
			sq.Or{
				sq.Eq{"id": ID},
				sq.Eq{"slug": ID},
			},
			sq.Eq{
				"deleted_time": nil,
			},
		}).ToSql()
	if err != nil {
		return models.Article{}, err
	}

	rows := r.db.QueryRowContext(ctx, query, args...)

	article := models.Article{}
	err = rows.Scan(
		&article.ID,
		&article.Title,
		&article.Slug,
		&article.Content,
		&article.CreateTime,
		&article.UpdateTime,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Article{}, gower.NotFoundErrorf("article not found")
		}
		return models.Article{}, err
	}

	return article, nil
}

func (r mysqlRepository) Store(ctx context.Context, article models.Article) error {
	qStr, args, err := sq.Insert("articles").
		Columns("id", "title", "slug", "content", "created_time", "updated_time").
		Values(article.ID, article.Title, article.Slug, article.Content, article.CreateTime, article.UpdateTime).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, qStr, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r mysqlRepository) Update(ctx context.Context, article models.Article) error {
	qStr, args, err := sq.Update("articles").
		SetMap(sq.Eq{
			"title":        article.Title,
			"content":      article.Content,
			"updated_time": article.UpdateTime,
		}).
		Where(sq.Eq{"id": article.ID}).Where("deleted_time IS NULL").ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.ExecContext(ctx, qStr, args...)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()

	if affected != 1 {
		return gower.NotFoundErrorf("article not found")
	}

	return nil
}

func (r mysqlRepository) Delete(ctx context.Context, ID string) error {
	qStr, args, err := sq.Update("articles").
		SetMap(sq.Eq{
			"deleted_time": time.Now(),
		}).
		Where(sq.Eq{"id": ID}).Where("deleted_time IS NULL").ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.ExecContext(ctx, qStr, args...)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()

	if affected != 1 {
		return gower.NotFoundErrorf("article not found")
	}

	return nil
}
