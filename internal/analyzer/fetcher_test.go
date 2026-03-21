package analyzer

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestFetchURL_Success tests that fetchURL returns a response for a valid URL
func TestFetchURL_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("<html><body>Hi</body></html>")); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	resp, err := fetchURL(server.URL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

// TestFetchURL_NotFound tests that fetchURL returns an error for a 404 response
func TestFetchURL_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := fetchURL(server.URL)
	if err == nil {
		t.Fatal("expected an error for 404 response, got nil")
	}
}

// TestFetchURL_Unreachable tests that fetchURL returns an error for an unreachable URL
func TestFetchURL_Unreachable(t *testing.T) {
	_, err := fetchURL("http://localhost:19999")
	if err == nil {
		t.Fatal("expected an error for unreachable URL, got nil")
	}
}

// TestFetchURL_ReturnsStatusCode tests that FetchError carries the correct status code
func TestFetchURL_StatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := fetchURL(server.URL)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var fetchErr *FetchError
	if !errors.As(err, &fetchErr) {
		t.Fatalf("expected FetchError, got %T", err)
	}

	if fetchErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", fetchErr.StatusCode)
	}
}
