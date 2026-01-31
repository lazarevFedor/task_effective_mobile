// Package main contains the executable entry point for the subscriptions service.
//
// This package configures a structured JSON logger, creates a context that
// carries the logger, and starts the HTTP server implemented in the
// internal/server package. Any fatal error returned by the server is logged
// and causes the process to exit with a non-zero status code.
package main

import (
	"context"
	"log/slog"
	"os"
	"task_effective_mobile/internal/server"
	"task_effective_mobile/pkg/logger"
)

// main configures structured logging and starts the HTTP server.
//
// The function constructs a context containing a JSON logger and calls
// server.Start. If server.Start returns an error, it is logged and the
// process exits with status code 1.
func main() {
	ctx := logger.WithLogger(context.Background(), slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	if err := server.Start(ctx); err != nil {
		logger.GetLogger(ctx).Error("server exited with error", "err", err)
		os.Exit(1)
	}
}
