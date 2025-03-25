package request

import (
	"net/http"

	"github.com/dexlabsio/garlic/pkg/logging"
	"go.uber.org/zap"
)

const (
	RequestIdHeaderKey = "X-Request-ID"
	SessionIdHeaderKey = "X-Session-ID"
)

// GetLogger returns the logger from the given request's context
func GetLogger(r *http.Request) *zap.Logger {
	ctx := r.Context()
	return logging.GetLoggerFromContext(ctx)
}

// GetRequestId is a helper function that retrieves the request ID from a request
func GetRequestId(r *http.Request) string {
	requestId := r.Header.Get(RequestIdHeaderKey)
	if requestId == "" {
		return "request-id-not-set"
	}

	return requestId
}

// GetSessionId is a helper function that retrieves the session ID from a request
func GetSessionId(r *http.Request) string {
	sessionId := r.Header.Get(SessionIdHeaderKey)
	if sessionId == "" {
		return "session-id-not-set"
	}

	return sessionId
}
