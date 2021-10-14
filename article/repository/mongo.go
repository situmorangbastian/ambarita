package repository

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/situmorangbastian/ambarita/models"
	"github.com/situmorangbastian/eclipse"
)

const (
	collectionArticle = "article"
)

type mongoRepository struct {
	db *mongo.Database
}

// NewMongoRepository will create the article mongo repository
func NewMongoRepository(db *mongo.Database) models.ArticleRepository {
	return mongoRepository{
		db: db,
	}
}

func (r mongoRepository) Fetch(ctx context.Context, cursor string, num int) ([]models.Article, string, error) {
	mongoOpts := options.Find().
		SetLimit(int64(num)).
		SetProjection(bson.M{
			"_id":          true,
			"slug":         true,
			"title":        true,
			"content":      true,
			"created_time": true,
			"updated_time": true,
		})
	mongoOpts.SetSort(bson.M{"created_time": -1})

	mongoFilter := bson.D{}

	if cursor != "" {
		nameFilter := bson.M{}
		name, err := DecodeCursor(cursor)
		if err != nil {
			return []models.Article{}, "", err
		}
		nameFilter["$lt"] = name

		if len(nameFilter) > 0 {
			mongoFilter = append(mongoFilter, bson.E{Key: "created_time", Value: nameFilter})
		}
	}

	mongoFilter = append(mongoFilter, bson.E{Key: "deleted_time", Value: bson.M{"$exists": false}})

	cur, err := r.db.Collection(collectionArticle).Find(ctx, mongoFilter, mongoOpts)
	if err != nil {
		return []models.Article{}, "", errors.Wrap(err, "find the documents")
	}

	results := make([]models.Article, 0)
	if err := cur.All(ctx, &results); err != nil {
		return results, "", errors.Wrap(err, "decodes each document")
	}

	if err := cur.Err(); err != nil {
		return results, "", errors.Wrap(err, "processing the cursor")
	}

	nextCursor := cursor
	if len(results) > 0 {
		nextCursor = EncodeCursor(results[len(results)-1].CreateTime)
	}

	return results, nextCursor, nil
}

func (r mongoRepository) Get(ctx context.Context, ID string) (models.Article, error) {
	res := r.db.Collection(collectionArticle).FindOne(
		ctx,
		bson.D{
			bson.E{
				Key: "$or",
				Value: []bson.D{
					{bson.E{Key: "_id", Value: ID}},
					{bson.E{Key: "slug", Value: ID}},
				},
			},
			bson.E{
				Key: "deleted_time",
				Value: bson.D{
					bson.E{Key: "$exists", Value: false},
				},
			},
		},
		options.FindOne().SetProjection(bson.M{
			"_id":          true,
			"slug":         true,
			"title":        true,
			"content":      true,
			"created_time": true,
			"updated_time": true,
		}),
	)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return models.Article{}, eclipse.NotFoundErrorf("article not found")
		}
		return models.Article{}, res.Err()
	}

	var article models.Article
	if err := res.Decode(&article); err != nil {
		return models.Article{}, err
	}
	return article, nil
}

func (r mongoRepository) Store(ctx context.Context, article models.Article) error {
	if _, err := r.db.Collection(collectionArticle).InsertOne(ctx, article); err != nil {
		return err
	}

	return nil
}

func (r mongoRepository) Update(ctx context.Context, article models.Article) error {
	res, err := r.db.Collection(collectionArticle).UpdateOne(ctx, bson.M{"_id": article.ID}, bson.M{
		"$set": bson.M{
			"title":        article.Title,
			"content":      article.Content,
			"updated_time": article.UpdateTime,
		},
	})

	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return err
	}

	return nil
}

func (r mongoRepository) Delete(ctx context.Context, ID string) error {
	res, err := r.db.Collection(collectionArticle).UpdateOne(ctx, bson.M{"_id": ID}, bson.M{
		"$set": bson.M{
			"deleted_time": time.Now(),
		},
	})

	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return err
	}

	return nil
}
