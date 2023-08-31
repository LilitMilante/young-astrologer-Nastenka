package api

import (
	"encoding/json"
	"net/http"
	"time"

	"young-astrologer-Nastenka/internal/service"
)

type Handler struct {
	s *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Image(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")

	if date == "" {
		date = time.Now().UTC().Format("2006-01-02")
	} else {
		_, err := time.Parse("2006-01-02", date)
		if err != nil {
			http.Error(w, "date must be in YYYY-MM-DD format", http.StatusBadRequest)
			return
		}
	}

	apod, err := h.s.APODByDate(r.Context(), date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(apod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Images(w http.ResponseWriter, r *http.Request) {
	apods, err := h.s.APODs(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(apods)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
