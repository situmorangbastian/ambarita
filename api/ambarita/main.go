package ambarita

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/situmorangbastian/ambarita/api/vercelhttp"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	router.HandleFunc("/api/articles", vercelhttp.FetchAllArticles).Methods(http.MethodGet)

	router.ServeHTTP(w, r)
}
