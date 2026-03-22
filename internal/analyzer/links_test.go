package analyzer

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestIsInternalLink tests that internal and external links are classified correctly
func TestIsInternalLink_Internal(t *testing.T) {
	baseURL := "https://example.com"

	tests := []struct {
		link     string
		expected bool
	}{
		{"/about", true},                    // relative URL
		{"https://example.com/about", true}, // same host
		{"https://google.com", false},       // different host
		{"https://sub.example.com", false},  // subdomain is external
	}

	for _, tt := range tests {
		result := isInternalLink(tt.link, baseURL)
		if result != tt.expected {
			t.Errorf("isInternalLink(%q, %q) = %v, expected %v", tt.link, baseURL, result, tt.expected)
		}
	}
}

// TestResolveLink tests that links are resolved correctly
func TestResolveLink(t *testing.T) {
	baseURL := "https://example.com"

	tests := []struct {
		link     string
		expected string
	}{
		{"#section", ""},                        // anchor -> ignored
		{"", ""},                                // empty -> ignored
		{"/about", "https://example.com/about"}, // relative -> resolved
		{"https://google.com", "https://google.com"}, // absolute -> unchanged
	}

	for _, tt := range tests {
		result := resolveLink(tt.link, baseURL)
		if result != tt.expected {
			t.Errorf("resolveLink(%q, %q) = %q, expected %q", tt.link, baseURL, result, tt.expected)
		}
	}
}

// TestIsLinkAccessible_Accessible tests that accessible links return true
func TestIsLinkAccessible_Accessible(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	a := New()
	if !a.isLinkAccessible(context.Background(), server.URL) {
		t.Error("expected link to be accessible")
	}
}

// TestIsLinkAccessible_NotFound tests that 404 links return false
func TestIsLinkAccessible_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	a := New()
	if a.isLinkAccessible(context.Background(), server.URL) {
		t.Error("expected link to be inaccessible")
	}
}

// TestIsLinkAccessible_Unreachable tests that unreachable links return false
func TestIsLinkAccessible_Unreachable(t *testing.T) {
	a := New()
	if a.isLinkAccessible(context.Background(), "http://localhost:19999") {
		t.Error("expected link to be inaccessible")
	}
}

// TestCheckLinks tests that checkLinks returns correct summary
func TestCheckLinks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// build an external link using a different fake server
	externalServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer externalServer.Close()

	links := []string{
		externalServer.URL,    // external - different host
		server.URL + "/about", // internal - same host as baseURL
	}

	a := New()
	summary := a.checkLinks(context.Background(), links, server.URL)

	if summary.InternalCount != 1 {
		t.Errorf("expected 1 internal link, got %d", summary.InternalCount)
	}
	if summary.ExternalCount != 1 {
		t.Errorf("expected 1 external link, got %d", summary.ExternalCount)
	}
	if summary.InaccessibleCount != 0 {
		t.Errorf("expected 0 inaccessible links, got %d", summary.InaccessibleCount)
	}
}
