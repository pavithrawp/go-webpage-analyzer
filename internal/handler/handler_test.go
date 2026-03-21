package handler

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/pavithrawp/go-webpage-analyzer/internal/analyzer"
)

// setup creates a test handler with a real analyzer and logger
func setup() *Handler {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	a := analyzer.New()
	tmpl := template.Must(template.New("index.html").Parse(`<html></html>`))
	return New(logger, a, tmpl)
}

// TestIndex tests that the index handler returns 200
func TestIndex(t *testing.T) {
	h := setup()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.Index(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

// TestAnalyze_InvalidBody tests that an invalid request body returns 400
func TestAnalyze_InvalidBody(t *testing.T) {
	h := setup()

	req := httptest.NewRequest(http.MethodPost, "/analyze", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	h.Analyze(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp errorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Success {
		t.Error("expected success to be false")
	}
}

// TestAnalyze_EmptyURL tests that an empty URL returns 400
func TestAnalyze_EmptyURL(t *testing.T) {
	h := setup()

	body := `{"url": ""}`
	req := httptest.NewRequest(http.MethodPost, "/analyze", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.Analyze(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// TestAnalyze_InvalidURL tests that an invalid URL returns 400
func TestAnalyze_InvalidURL(t *testing.T) {
	h := setup()

	body := `{"url": "abcd"}`
	req := httptest.NewRequest(http.MethodPost, "/analyze", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.Analyze(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp errorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Success {
		t.Error("expected success to be false")
	}
}

// TestAnalyze_ValidURL tests that a valid URL returns 200 with a result
func TestAnalyze_ValidURL(t *testing.T) {
	// create a fake server to analyze
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`<!DOCTYPE html>
		<html>
			<head><title>Test Page</title></head>
			<body><h1>Hello</h1></body>
		</html>`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	h := setup()

	body, _ := json.Marshal(map[string]string{"url": server.URL})
	req := httptest.NewRequest(http.MethodPost, "/analyze", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Analyze(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result["title"] != "Test Page" {
		t.Errorf("expected title 'Test Page', got %v", result["title"])
	}
}

// TestAnalyze_UnreachableURL tests that an unreachable URL returns an error
func TestAnalyze_UnreachableURL(t *testing.T) {
	h := setup()

	body := `{"url": "http://localhost:19999"}`
	req := httptest.NewRequest(http.MethodPost, "/analyze", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.Analyze(w, req)

	if w.Code == http.StatusOK {
		t.Error("expected error status for unreachable URL")
	}

	var resp errorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Success {
		t.Error("expected success to be false")
	}
}
