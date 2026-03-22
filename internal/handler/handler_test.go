package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/pavithrawp/go-webpage-analyzer/internal/analyzer"
)

type mockAnalyzer struct {
	result *analyzer.Result
	err    error
}

func (m *mockAnalyzer) Analyze(url string) (*analyzer.Result, error) {
	return m.result, m.err
}

// setup creates a test handler with a real analyzer and logger
func setup(result *analyzer.Result, err error) *Handler {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	tmpl := template.Must(template.New("index.html").Parse(`<html></html>`))
	return New(logger, &mockAnalyzer{result: result, err: err}, tmpl)
}

// TestIndex tests that the index handler returns 200
func TestIndex(t *testing.T) {
	h := setup(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.Index(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

// TestAnalyze_InvalidBody tests that an invalid request body returns 400
func TestAnalyze_InvalidBody(t *testing.T) {
	h := setup(nil, nil)

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
	h := setup(nil, nil)

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
	h := setup(nil, nil)

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
	res := &analyzer.Result{
		HTMLVersion:  "HTML5",
		Title:        "Test Page",
		Headings:     map[string]int{"h1": 1},
		HasLoginForm: false,
	}

	h := setup(res, nil)

	body, _ := json.Marshal(map[string]string{"url": "https://abcd.com"})
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
	h := setup(nil, fmt.Errorf("something went wrong"))

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
