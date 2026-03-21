package main

import (
	"log/slog"
	"net/http"
	"os"

	_ "net/http/pprof"

	"github.com/joho/godotenv"
	"github.com/pavithrawp/go-webpage-analyzer/internal/analyzer"
	"github.com/pavithrawp/go-webpage-analyzer/internal/handler"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		logger.Warn("no .env file found, using system environment")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if not set
	}

	mux := http.NewServeMux()

	pageAnalyzer := analyzer.New()
	// register routes
	h := handler.New(logger, pageAnalyzer)
	mux.HandleFunc("GET /{$}", h.Index)
	mux.HandleFunc("POST /analyze", h.Analyze)

	// register pprof routes
	mux.HandleFunc("/debug/pprof/", handler.PprofAuth(http.DefaultServeMux.ServeHTTP))
	mux.HandleFunc("/debug/pprof/cmdline", handler.PprofAuth(http.DefaultServeMux.ServeHTTP))
	mux.HandleFunc("/debug/pprof/profile", handler.PprofAuth(http.DefaultServeMux.ServeHTTP))
	mux.HandleFunc("/debug/pprof/symbol", handler.PprofAuth(http.DefaultServeMux.ServeHTTP))
	mux.HandleFunc("/debug/pprof/trace", handler.PprofAuth(http.DefaultServeMux.ServeHTTP))

	logger.Info("server starting", "port", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		logger.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
