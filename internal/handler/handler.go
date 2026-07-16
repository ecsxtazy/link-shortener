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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "invalid request body",
		})
		return
	}
	short, err := h.service.Shorten(r.Context(), req.URL)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "internal server error",
		})
		return
	}
	if req.URL == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "invalid request body",
		})
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "link not found"})
			return
		}
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)

			json.NewEncoder(w).Encode(ErrorResponse{
				Error: "internal server error",
			})
			return
		}
	}
	resp := ResolveResponse{
		URL: url,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
