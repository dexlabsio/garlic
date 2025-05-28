package middleware

import (
	"context"
	"net/http"
)

// ContextCancel is a middleware that wraps an HTTP handler to provide a context
// with a cancellation function. This allows the context to be cancelled when
// the request is complete, helping to manage resources and prevent leaks.
// It creates a new context with a cancel function, attaches it to the request,
// and ensures the cancel function is called after the request is processed.
func ContextCancel(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
