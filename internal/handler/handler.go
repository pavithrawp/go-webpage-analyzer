package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/pavithrawp/go-webpage-analyzer/internal/analyzer"
)

type Handler struct {
	logger   *slog.Logger
	analyzer *analyzer.Analyzer
}

func New(logger *slog.Logger, a *analyzer.Analyzer) *Handler {
	return &Handler{
		logger:   logger,
		analyzer: a,
	}
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
	var req analyzeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	result, err := h.analyzer.Analyze(req.URL)
	if err != nil {
		h.logger.Error("failed to analyze URL", "url", req.URL, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}

}
