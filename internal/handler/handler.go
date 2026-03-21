package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"html/template"

	"github.com/pavithrawp/go-webpage-analyzer/internal/analyzer"
	"github.com/pavithrawp/go-webpage-analyzer/internal/validator"
)

const (
	contentTypeJSON   = "application/json"
	headerContentType = "Content-Type"
)

type Handler struct {
	logger   *slog.Logger
	analyzer *analyzer.Analyzer
	template *template.Template
}

func New(logger *slog.Logger, pageAnalyzer *analyzer.Analyzer, tmpl *template.Template) *Handler {
	return &Handler{
		logger:   logger,
		analyzer: pageAnalyzer,
		template: tmpl,
	}
}

// Index handles GET / and serves the analyzer form
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if err := h.template.ExecuteTemplate(w, "index.html", nil); err != nil {
		h.logger.Error("failed to render template", "error", err)
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Analyze handles POST /analyze and processes the URL analysis request
func (h *Handler) Analyze(w http.ResponseWriter, r *http.Request) {
	var req analyzeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.ValidateURL(req.URL); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.analyzer.Analyze(req.URL)
	if err != nil {
		h.logger.Error("failed to analyze URL", "url", req.URL, "error", err)

		// check if it's a FetchError to return the correct HTTP status code
		var fetchErr *analyzer.FetchError
		if errors.As(err, &fetchErr) {
			h.writeError(w, fetchErr.StatusCode, fetchErr.Error())
			return
		}

		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set(headerContentType, contentTypeJSON)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}

}

// writeError writes a JSON error response with the given status code and message
func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set(headerContentType, contentTypeJSON)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(errorResponse{
		Success: false,
		Error:   message,
	}); err != nil {
		h.logger.Error("failed to encode error response", "error", err)
	}
}
