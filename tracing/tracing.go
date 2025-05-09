package tracing

import (
	"context"

	"github.com/dexlabsio/garlic/errors"
	"github.com/dexlabsio/garlic/logging"

	"github.com/google/uuid"
)

type key int

const (
	RequestIdKey key = iota
	SessionIdKey
)

// GetRequestIdFromContext is a helper function that retrieves the request ID from a context
func GetRequestIdFromContext(ctx context.Context) (uuid.UUID, error) {
	val := ctx.Value(RequestIdKey)
	if val == nil {
		return uuid.Nil, errors.New(
			KindContextValueNotFoundError,
			"request id is not set in this context",
		)
	}

	requestId, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New(
			KindContextError,
			"invalid request id found in context",
			errors.Context(
				errors.Field("invalid_request_id", val),
			),
		)
	}

	return requestId, nil
}

// SetContextRequestId is a helper function that associates a request ID with a context
// by storing the request ID in the context using a predefined key. This allows
// the request ID to be retrieved later from the context, enabling consistent
// tracking of request information throughout the request lifecycle.
func SetContextRequestId(ctx context.Context, requestId uuid.UUID) context.Context {
	return context.WithValue(ctx, RequestIdKey, requestId)
}

// GetSessionIdFromContext is a helper function that retrieves the session ID from a context
func GetSessionIdFromContext(ctx context.Context) (string, error) {
	val := ctx.Value(SessionIdKey)
	if val == nil {
		return "", errors.New(
			KindContextValueNotFoundError,
			"session id is not set in this context",
		)
	}

	sessionId, ok := val.(string)
	if !ok {
		return "", errors.New(
			KindContextError,
			"invalid session id found in context",
			errors.Context(
				errors.Field("invalid_session_id", val),
			),
		)
	}

	return sessionId, nil
}

// MustGetRequestIdFromContext is a utility function that retrieves the request ID from the context.
// It panics if the request ID is not found or if there is an error during retrieval.
// This function is useful in scenarios where the presence of a request ID is critical,
// and the application cannot proceed without it. It logs a fatal error message before
// terminating the application if the request ID is missing or invalid.
func MustGetRequestIdFromContext(ctx context.Context) uuid.UUID {
	requestId, err := GetRequestIdFromContext(ctx)
	if err != nil {
		err = errors.Propagate(err, "Failed to get request ID from context")
		logging.Global().Fatal("Failed to get request ID from context", errors.Zap(err))
	}

	return requestId
}

// MustGetSessionIdFromContext is a utility function that retrieves the session ID from the context.
// It panics if the session ID is not found or if there is an error during retrieval.
// This function is useful in scenarios where the presence of a session ID is critical,
// and the application cannot proceed without it. It logs a fatal error message before
// terminating the application if the session ID is missing or invalid.
func MustGetSessionIdFromContext(ctx context.Context) string {
	sessionId, err := GetSessionIdFromContext(ctx)
	if err != nil {
		err = errors.Propagate(err, "Failed to get session ID from context")
		logging.Global().Fatal("Failed to get session ID from context", errors.Zap(err))
	}

	return sessionId
}

// SetContextSessionId is a helper function that associates a session ID with a context
// by storing the session ID in the context using a predefined key. This allows
// the session ID to be retrieved later from the context, enabling consistent
// tracking of session information throughout the request lifecycle.
func SetContextSessionId(ctx context.Context, sessionId string) context.Context {
	return context.WithValue(ctx, SessionIdKey, sessionId)
}
