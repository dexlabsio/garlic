package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dexlabsio/garlic/pkg/logging"
	"github.com/dexlabsio/garlic/pkg/request"
	"go.uber.org/zap"
)

// The NewLoggingMiddleware function generates a unique request ID for each
// incoming HTTP request, enriches the logger with this ID, and stores both
// the logger and the ID in the request's context for future use in
// subsequent layers.
func NewLoggingMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If the request is a health check, we don't need to log it.
			if r.URL.String() == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()

			// The request context coming in must already contain the request and session IDs.
			// We will keep the same request ID in the context to maintain traceability.
			requestId := request.GetRequestId(r)
			sessionId := request.GetSessionId(r)

			lrw := &loggingResponseWriter{w, http.StatusOK, 0}
			lrw.Header().Set(request.RequestIdHeaderKey, requestId)
			lrw.Header().Set(request.SessionIdHeaderKey, sessionId)

			logger := logger.With(
				zap.String("request_id", requestId),
				zap.String("session_id", sessionId),
				zap.String("request_method", r.Method),
				zap.String("request_url", r.URL.String()),
			)
			logger.Debug(fmt.Sprintf("Started handling %s request for %s", r.Method, r.URL.String()))

			// Set the logger in the context for future use
			ctx := r.Context()
			ctx = context.WithValue(ctx, logging.LoggerKey, logger)
			r = r.WithContext(ctx)

			next.ServeHTTP(lrw, r)

			duration := time.Since(start)

			logger = logger.With(
				zap.Int("response_status", lrw.statusCode),
				zap.Duration("response_time", duration),
				zap.Int("response_size", lrw.responseSize),
			)

			logger.Info(fmt.Sprintf(
				"Response Status: %d | URI: %s %s | Duration: %s | Response Size: %d bytes",
				lrw.statusCode,
				r.Method,
				r.URL.String(),
				duration,
				lrw.responseSize,
			))
		})
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter

	statusCode   int
	responseSize int
}

// WriteHeader writes the status code to the response writer and stores it in.
func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

// Write writes the data to the response writer and stores the size of the data.
func (w *loggingResponseWriter) Write(data []byte) (int, error) {
	size, err := w.ResponseWriter.Write(data)
	w.responseSize += size
	return size, err
}
