package middleware

import (
	"net/http"

	"github.com/dexlabsio/garlic/errors"
	"github.com/dexlabsio/garlic/request"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	RequestIdHeaderKey = "X-Request-ID"
	SessionIdHeaderKey = "X-Session-ID"
)

func Tracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := request.GetLogger(r)

		requestId := uuid.New()
		w.Header().Set(RequestIdHeaderKey, requestId.String())
		l = l.With(zap.Stringer("request_id", requestId))
		r = request.SetLogger(r, l)
		r = request.SetRequestId(r, requestId)

		next.ServeHTTP(w, r)
	})
}

func PropagateTracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := request.GetLogger(r)

		requestId, err := requestHeaderRequestId(r)
		if err != nil {
			l.Error("Failed to propagate header", errors.Zap(err))
		} else {
			w.Header().Set(RequestIdHeaderKey, requestId.String())
			l = l.With(zap.Stringer("request_id", requestId))
			r = request.SetLogger(r, l)
			r = request.SetRequestId(r, requestId)
		}

		sessionId, err := requestHeaderSessionId(r)
		if err != nil {
			l.Error("Failed to propagate header", errors.Zap(err))
		} else {
			w.Header().Set(SessionIdHeaderKey, sessionId)
			l = l.With(zap.String("session_id", sessionId))
			r = request.SetLogger(r, l)
			r = request.SetSessionId(r, sessionId)
		}

		next.ServeHTTP(w, r)
	})
}

// requestHeaderRequestId is a helper function that retrieves the request ID from a request
func requestHeaderRequestId(r *http.Request) (uuid.UUID, error) {
	requestId, err := uuid.Parse(r.Header.Get(RequestIdHeaderKey))
	if err != nil {
		return uuid.Nil, errors.PropagateAs(
			errors.KindInvalidRequestError,
			err,
			"missing mandatory request header",
			errors.Hint("Please provide mandatory request header: %s", RequestIdHeaderKey),
		)
	}

	return requestId, nil
}

// requestHeaderSessionId is a helper function that retrieves the session ID from a request
func requestHeaderSessionId(r *http.Request) (string, error) {
	sessionId := r.Header.Get(SessionIdHeaderKey)
	if sessionId == "" {
		return "", errors.New(
			errors.KindInvalidRequestError,
			"missing mandatory request header",
			errors.Hint("Please provide mandatory request header: %s", SessionIdHeaderKey),
		)
	}

	return sessionId, nil
}
