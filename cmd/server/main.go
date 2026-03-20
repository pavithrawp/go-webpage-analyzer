package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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

	// register routes
	h := handler.New(logger)
	mux.HandleFunc("GET /", h.Index)
	mux.HandleFunc("POST /analyze", h.Analyze)

	logger.Info("server starting", "port", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		logger.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
