package analyzer

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestAnalyze_Success tests the full analysis with a fake server
func TestAnalyze_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`<!DOCTYPE html>
		<html>
			<head><title>Test Page</title></head>
			<body>
				<h1>Hello</h1>
				<h2>World</h2>
				<a href="/internal">Internal</a>
				<a href="https://external.com">External</a>
			</body>
		</html>`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	a := New()
	result, err := a.Analyze(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.HTMLVersion != "HTML5" {
		t.Errorf("expected HTML5, got %s", result.HTMLVersion)
	}

	if result.Title != "Test Page" {
		t.Errorf("expected 'Test Page', got %s", result.Title)
	}

	if result.Headings["h1"] != 1 {
		t.Errorf("expected 1 h1, got %d", result.Headings["h1"])
	}

	if result.Headings["h2"] != 1 {
		t.Errorf("expected 1 h2, got %d", result.Headings["h2"])
	}
}

// TestAnalyze_UnreachableURL tests that an unreachable URL returns an error
func TestAnalyze_UnreachableURL(t *testing.T) {
	a := New()
	_, err := a.Analyze(context.Background(), "http://localhost:19999")
	if err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}

// TestAnalyze_NonHTMLContent tests that non-HTML content returns an error
func TestAnalyze_NonHTMLContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"key": "value"}`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	a := New()
	_, err := a.Analyze(context.Background(), server.URL)
	if err == nil {
		t.Fatal("expected error for non-HTML content, got nil")
	}
}

// TestAnalyze_LoginForm tests that login form detection works
func TestAnalyze_LoginForm(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`<!DOCTYPE html>
		<html>
			<body>
				<form>
					<input type="text" name="username"/>
					<input type="password" name="password"/>
				</form>
			</body>
		</html>`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	a := New()
	result, err := a.Analyze(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !result.HasLoginForm {
		t.Error("expected login form to be detected")
	}
}

// TestAnalyze_404URL tests that a 404 URL returns a FetchError
func TestAnalyze_404URL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	a := New()
	_, err := a.Analyze(context.Background(), server.URL)
	if err == nil {
		t.Fatal("expected error for 404 URL, got nil")
	}
}
