package middleware

import (
	"net/http"

	"github.com/dexlabsio/garlic/request"
	"github.com/google/uuid"
)

const (
	RequestIdHeaderKey = "X-Request-ID"
)

func Tracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The request context coming in must already contain the request and session IDs.
		// We will keep the same request ID in the context to maintain traceability.
		requestId := uuid.New()
		r = request.SetRequestId(r, requestId)
		w.Header().Set(RequestIdHeaderKey, requestId.String())

		next.ServeHTTP(w, r)
	})
}
