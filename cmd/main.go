package main

import (
	"context"
	"log/slog"
	"os"
	"task_effective_mobile/internal/server"
	"task_effective_mobile/pkg/logger"
)

func main() {
	ctx := logger.WithLogger(context.Background(), slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	if err := server.Start(ctx); err != nil {
		logger.GetLogger(ctx).Error("server exited with error", "err", err)
		os.Exit(1)
	}
}
