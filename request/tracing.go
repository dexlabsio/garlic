package request

import (
	"net/http"

	"github.com/dexlabsio/garlic/errors"
	"github.com/dexlabsio/garlic/tracing"
	"github.com/google/uuid"
)

// GetRequestId is a helper function that retrieves the request ID from a request
func GetRequestId(r *http.Request) (uuid.UUID, error) {
	ctx := r.Context()
	id, err := tracing.GetRequestIdFromContext(ctx)
	if err != nil {
		return uuid.Nil, errors.Propagate(err, "failed to get request id in this request")
	}

	return id, nil
}

// SetRequestId is a helper function that associates a request ID with an HTTP request
// by storing the request ID in the request's context. This allows the request ID
// to be consistently accessed throughout the request's lifecycle, facilitating
// traceability and logging of request-specific information.
func SetRequestId(r *http.Request, requestId uuid.UUID) *http.Request {
	ctx := tracing.SetContextRequestId(r.Context(), requestId)
	return r.WithContext(ctx)
}

// GetSessionId is a helper function that retrieves the session ID from a request
func GetSessionId(r *http.Request) (string, error) {
	ctx := r.Context()
	id, err := tracing.GetSessionIdFromContext(ctx)
	if err != nil {
		return "", errors.Propagate(err, "failed to get session id in this request")
	}

	return id, nil
}

// SetSessionId is a helper function that associates a session ID with an HTTP request
// by storing the session ID in the request's context. This allows the session ID
// to be consistently accessed throughout the request's lifecycle, facilitating
// traceability and logging of session-specific information.
func SetSessionId(r *http.Request, sessionId string) *http.Request {
	ctx := tracing.SetContextSessionId(r.Context(), sessionId)
	return r.WithContext(ctx)
}
