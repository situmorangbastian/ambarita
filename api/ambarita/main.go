package ambarita

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/situmorangbastian/ambarita/api/vercelhttp"
	articleRepository "github.com/situmorangbastian/ambarita/article/repository"
	articleUsecase "github.com/situmorangbastian/ambarita/article/usecase"
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

	router := mux.NewRouter()

	router.HandleFunc("/api/posts", vercelhttp.FetchAllArticles).Methods(http.MethodGet)
	router.HandleFunc("/api/posts/{id}", vercelhttp.GetByID).Methods(http.MethodGet)

	router.ServeHTTP(w, r)
}
