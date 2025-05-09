package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dexlabsio/garlic/errors"
	"github.com/dexlabsio/garlic/logging"
	"github.com/dexlabsio/garlic/request"
	"go.uber.org/zap"
)

// Logging function generates a unique request ID for each
// incoming HTTP request, enriches the logger with this ID, and stores both
// the logger and the ID in the request's context for future use in
// subsequent layers.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the request is a health check, we don't need to log it.
		if r.URL.String() == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		logger := logging.Global()
		start := time.Now()

		// The request context coming in must already contain the request and session IDs.
		// We will keep the same request ID in the context to maintain traceability.
		requestId, err := request.GetRequestId(r)
		if err != nil {
			logger.Debug("Request ID will not be logged for this request", errors.Zap(err))
		} else {
			logger = logger.With(zap.Stringer("request_id", requestId))
		}

		logger = logger.With(
			zap.String("request_method", r.Method),
			zap.String("request_url", r.URL.String()),
		)

		r = request.SetLogger(r, logger)

		lrw := &loggingResponseWriter{w, http.StatusOK, 0}
		logger.Debug(fmt.Sprintf("Handling %s %s", r.Method, r.URL.String()))
		next.ServeHTTP(lrw, r)

		duration := time.Since(start)

		logger = logger.With(
			zap.Int("response_status", lrw.statusCode),
			zap.Duration("response_time", duration),
			zap.Int("response_size", lrw.responseSize),
		)

		logger.Info(fmt.Sprintf(
			"[%d] %s %s",
			lrw.statusCode,
			r.Method,
			r.URL.String(),
		))
	})
}

// loggingResponseWriter is a custom HTTP response writer that captures
// the status code and the size of the response. It embeds the standard
// http.ResponseWriter and overrides the WriteHeader and Write methods to
// store the status code and accumulate the size of the response body.
// This allows for enhanced logging of HTTP responses, including the
// status code and the total size of the response sent to the client.
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
