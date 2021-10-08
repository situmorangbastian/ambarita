package ambarita

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/situmorangbastian/ambarita/api/vercelhttp"
	articleHandler "github.com/situmorangbastian/ambarita/article/http"
	articleRepository "github.com/situmorangbastian/ambarita/article/repository"
	articleUsecase "github.com/situmorangbastian/ambarita/article/usecase"
	"github.com/situmorangbastian/eclipse"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	mongoURI := os.Getenv("MONGO_URI")

	mongoClient, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(mongoURI).
			SetConnectTimeout(2*time.Second).
			SetServerSelectionTimeout(3*time.Second),
	)
	if err != nil {
		log.Fatal("Mongo connection failed: ", err.Error())
	}

	mongoDBName := os.Getenv("MONGO_DB_NAME")

	ar := articleRepository.NewMongoRepository(mongoClient.Database(mongoDBName))
	vercelhttp.ArticleUsecase = articleUsecase.NewArticleUsecase(ar)

	e := echo.New()
	e.Use(eclipse.Error())

	// Domain
	au := articleUsecase.NewArticleUsecase(ar)
	articleHandler.NewHandler(e, au)

	e.ServeHTTP(w, r)
}
