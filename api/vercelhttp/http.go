package vercelhttp

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/situmorangbastian/ambarita/models"
)

func FetchAllArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp, _ := json.Marshal(models.Posts)

	_, _ = w.Write(resp)
}

func GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var resp []byte

	for _, post := range models.Posts {
		if post.ID == id {
			resp, _ = json.Marshal(post)
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(resp)
}
