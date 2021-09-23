package handler

import (
	"net/http"

	"github.com/go-chi/chi"
)

func NewHealthCheckRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/status", StatusActive)
	return r
}

func StatusActive(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(`{"healt-status": "OK"}`))
	w.WriteHeader(http.StatusOK)
}
