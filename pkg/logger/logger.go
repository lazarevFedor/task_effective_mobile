// Package logger provides helpers to store and retrieve a structured slog.Logger
// in a context.Context value. Storing the logger in the context allows passing
// the logger through call chains without modifying many function signatures.
package logger

import (
	"context"
	"log/slog"
)

// key is an unexported type used to avoid collisions in context keys.
// Using a private type guarantees that other packages cannot accidentally
// use the same key value.
type key string

// loggerKey is the context key under which a *slog.Logger is stored.
const loggerKey = key("logger")

// WithLogger returns a new context that carries the provided logger value.
//
// Use this function at the top-level of request handling (for example in
// main) to attach a configured *slog.Logger to the context. Downstream code
// can obtain the logger via GetLogger.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// GetLogger retrieves the *slog.Logger stored in ctx using WithLogger.
//
// The function performs a type assertion and returns the logger pointer. If
// no logger is present in the context the type assertion will panic; callers
// should ensure that WithLogger was used to populate the context before
// calling GetLogger (for example by establishing the logger in main).
func GetLogger(ctx context.Context) *slog.Logger {
	logger := ctx.Value(loggerKey)
	return logger.(*slog.Logger)
}
