package router

import (
	"link-shortener/internal/handler"
	"net/http"
)

func New(handler *handler.Handler) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /shorten", handler.Shorten)

	mux.HandleFunc("GET /{code}", handler.Resolve)

	return mux
}
