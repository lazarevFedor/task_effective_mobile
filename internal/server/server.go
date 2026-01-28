package server

import (
	"context"
	"fmt"
	"net/http"
	"task_effective_mobile/internal/config"
)

func Start(ctx context.Context) error {
	mux := http.NewServeMux()
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("start: failed to load config: %w", err)
	}
	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), mux)
	if err != nil {
		return fmt.Errorf("start: error while starting http server: %w", err)
	}
	return nil
}
