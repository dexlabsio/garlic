package logging

import (
	"context"

	"go.uber.org/zap"
)

type key int

const (
	LoggerKey key = iota
)

// GetLoggerFromContext is a helper function that retrieves the logger from a context
func GetLoggerFromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(LoggerKey).(*zap.Logger)
	if !ok {
		panic("Failed to get logger from context")
	}

	return logger
}

// SetContextLogger is a helper function that associates a logger with a context
// by storing the logger in the context using a predefined key. This allows
// the logger to be retrieved later from the context, enabling consistent
// logging throughout the request lifecycle.
func SetContextLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}
