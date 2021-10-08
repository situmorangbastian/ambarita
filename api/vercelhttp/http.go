package vercelhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/situmorangbastian/ambarita/models"
	"github.com/situmorangbastian/gower"
)

var (
	ArticleUsecase models.ArticleUsecase
)

func FetchAllArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	num := 20
	cursor := r.URL.Query().Get("cursor")

	if r.URL.Query().Get("num") != "" {
		numInt, err := strconv.Atoi(r.URL.Query().Get("num"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message": "invalid num parameter"}`))
			return
		}
		num = numInt
	}

	articles, nextCursor, err := ArticleUsecase.Fetch(context.Background(), cursor, num)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message": "internal server error"}`))
		return
	}

	resp, _ := json.Marshal(articles)

	w.Header().Add("X-Cursor", nextCursor)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	article, err := ArticleUsecase.Get(context.Background(), id)
	if err != nil {
		switch errors.Cause(err).(type) {
		case gower.NotFoundError:
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message": "post not found"}`))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message": "internal server error"}`))
			return
		}
	}

	resp, _ := json.Marshal(article)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}
