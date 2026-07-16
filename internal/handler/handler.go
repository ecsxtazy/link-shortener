package handler

import (
	"encoding/json"
	"errors"
	"link-shortener/internal/repository"
	"link-shortener/internal/service"
	"net/http"
)

type Handler struct {
	service service.Shortener
}

func New(service service.Shortener) *Handler {
	return &Handler{
		service: service,
	}
}
func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	short, err := h.service.Shorten(r.Context(), req.URL)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if req.URL == "" {
		handleError(w, http.StatusBadRequest, "url is required")
		return
	}
	resp := ShortenResponse{
		ShortURL: short,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Resolve(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	url, err := h.service.Resolve(r.Context(), code)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			handleError(w, http.StatusNotFound, err.Error())
			return
		}
		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := ResolveResponse{
		URL: url,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
