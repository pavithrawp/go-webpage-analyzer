package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestPprofAuth_NoCredentials tests that missing credentials return 401
func TestPprofAuth_NoCredentials(t *testing.T) {
	if err := os.Setenv(pprofUsernameEnv, "admin"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	if err := os.Setenv(pprofPasswordEnv, "admin"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	defer func() { _ = os.Unsetenv(pprofUsernameEnv) }()
	defer func() { _ = os.Unsetenv(pprofPasswordEnv) }()

	handler := PprofAuth(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/debug/pprof/", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

// TestPprofAuth_InvalidCredentials tests that wrong credentials return 401
func TestPprofAuth_InvalidCredentials(t *testing.T) {
	if err := os.Setenv(pprofUsernameEnv, "admin"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	if err := os.Setenv(pprofPasswordEnv, "admin"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	defer func() { _ = os.Unsetenv(pprofUsernameEnv) }()
	defer func() { _ = os.Unsetenv(pprofPasswordEnv) }()

	handler := PprofAuth(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/debug/pprof/", nil)
	req.SetBasicAuth("test", "test")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

// TestPprofAuth_ValidCredentials tests that correct credentials return 200
func TestPprofAuth_ValidCredentials(t *testing.T) {
	if err := os.Setenv(pprofUsernameEnv, "admin"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	if err := os.Setenv(pprofPasswordEnv, "admin"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	defer func() { _ = os.Unsetenv(pprofUsernameEnv) }()
	defer func() { _ = os.Unsetenv(pprofPasswordEnv) }()

	handler := PprofAuth(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/debug/pprof/", nil)
	req.SetBasicAuth("admin", "admin")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
