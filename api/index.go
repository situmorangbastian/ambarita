package handler

import (
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello from Go!</h1>"))
}
