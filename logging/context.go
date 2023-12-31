package logging

import (
	"context"
	"log/slog"
)

var contextKey = struct{}{}

// FromContext extracts a logger from the context.
// If there's no logger in the context, it will return
// a new default logger.
func FromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(contextKey).(*slog.Logger)
	if ok {
		return logger
	}

	return Default()
}

// ToContext returns the appended context with the saved logger in.
func ToContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKey, logger)
}
