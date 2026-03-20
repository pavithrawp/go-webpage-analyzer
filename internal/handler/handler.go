package handler

import (
	"log/slog"
	"net/http"
)

type Handler struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *Handler {
	return &Handler{logger: logger}
}

// Index handles GET / and returns the status of the server
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Webpage Analyzer is running")); err != nil {
		h.logger.Error("failed to write response", "error", err)
	}
}

// Analyze handles POST /analyze and processes the URL analysis request
func (h *Handler) Analyze(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Analyze endpoint hit")); err != nil {
		h.logger.Error("failed to write response", "error", err)
	}
}
