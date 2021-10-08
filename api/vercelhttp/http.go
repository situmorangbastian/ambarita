package vercelhttp

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/situmorangbastian/ambarita/models"
)

var (
	ArticleUsecase models.ArticleUsecase
)

func FetchAllArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	data, _, _ := ArticleUsecase.Fetch(context.Background(), "", 20)

	resp, _ := json.Marshal(data)

	_, _ = w.Write(resp)
}

func GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	data, _ := ArticleUsecase.Get(context.Background(), id)

	resp, _ := json.Marshal(data)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(resp)
}
