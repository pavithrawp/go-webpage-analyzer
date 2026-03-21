package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/joho/godotenv"
	"github.com/pavithrawp/go-webpage-analyzer/internal/analyzer"
	"github.com/pavithrawp/go-webpage-analyzer/internal/handler"
)

// shutdownTimeout is the maximum time to wait for the server to shutdown gracefully
const shutdownTimeout = 30 * time.Second

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		logger.Warn("no .env file found, using system environment")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default to port 8080 if not set
	}

	mux := http.NewServeMux()

	pageAnalyzer := analyzer.New()

	// register routes
	h := handler.New(logger, pageAnalyzer)
	mux.HandleFunc("GET /{$}", h.Index)
	mux.HandleFunc("POST /analyze", h.Analyze)

	// register pprof routes protected with basic auth
	mux.HandleFunc("/debug/pprof/", handler.PprofAuth(http.DefaultServeMux.ServeHTTP))
	mux.HandleFunc("/debug/pprof/cmdline", handler.PprofAuth(http.DefaultServeMux.ServeHTTP))
	mux.HandleFunc("/debug/pprof/profile", handler.PprofAuth(http.DefaultServeMux.ServeHTTP))
	mux.HandleFunc("/debug/pprof/symbol", handler.PprofAuth(http.DefaultServeMux.ServeHTTP))
	mux.HandleFunc("/debug/pprof/trace", handler.PprofAuth(http.DefaultServeMux.ServeHTTP))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// start server in a goroutine so it does not block
	go func() {
		logger.Info("server starting", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	// give active requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped")
}
