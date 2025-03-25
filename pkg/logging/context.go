package logging

import (
	"context"

	"go.uber.org/zap"
)

type key int

const (
	LoggerKey key = iota
	RequestIdKey
	SessionIdKey
)

// GetLoggerFromContext is a helper function that retrieves the logger from a context
func GetLoggerFromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(LoggerKey).(*zap.Logger)
	if !ok {
		panic("Failed to get logger from context")
	}

	return logger
}

// GetRequestIdFromContext is a helper function that retrieves the request ID from a context
func GetRequestIdFromContext(ctx context.Context) string {
	requestId, ok := ctx.Value(RequestIdKey).(string)
	if !ok {
		panic("Failed to get request ID from context")
	}

	return requestId
}

// GetSessionIdFromContext is a helper function that retrieves the session ID from a context
func GetSessionIdFromContext(ctx context.Context) string {
	sessionId, ok := ctx.Value(SessionIdKey).(string)
	if !ok {
		panic("Failed to get session ID from context")
	}

	return sessionId
}
