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

	articleHandler "github.com/situmorangbastian/ambarita/article/http"
	articleRepository "github.com/situmorangbastian/ambarita/article/repository"
	articleUsecase "github.com/situmorangbastian/ambarita/article/usecase"
	"github.com/situmorangbastian/ambarita/cmd/logger"
	"github.com/situmorangbastian/eclipse"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	logger.SetupLogs()

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

	e := echo.New()
	e.Use(eclipse.Error())
	g := e.Group("/api")

	// Domain
	au := articleUsecase.NewArticleUsecase(ar)
	articleHandler.NewGroupHandler(g, au)

	e.ServeHTTP(w, r)
}
