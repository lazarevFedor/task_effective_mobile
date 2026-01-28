package server

import (
	"context"
	"fmt"
	"net/http"
)

func Start(ctx context.Context) error {
	mux := http.NewServeMux()
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return fmt.Errorf("start: error while starting http server: %w", err)
	}
	return nil
}
