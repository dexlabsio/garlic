package middleware

import (
	"net/http"
	"time"

	"github.com/dexlabsio/garlic/pkg/monitoring"
	"github.com/go-chi/chi/v5"
)

// statusRecorder is a custom response writer that records the HTTP status code.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

// WriteHeader overrides the default WriteHeader method to record the status code.
func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

// MetricsMonitor is a middleware that monitors and records metrics for incoming requests.
func MetricsMonitor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := statusRecorder{w, 200}
		method := r.Method
		route := getRoutePattern(r)

		if isIgnoredRoute(route) {
			next.ServeHTTP(&rec, r)
			return
		}

		monitoring.IncrementActiveRequests(method, route)

		defer monitorLatency(method, route, time.Now())
		next.ServeHTTP(&rec, r)

		monitoring.DecrementActiveRequests(method, route)
		monitoring.IncrementTraffic(method, route, rec.status)
	})
}

// getRoutePattern retrieves the route pattern from the request context.
func getRoutePattern(r *http.Request) string {
	rctx := chi.RouteContext(r.Context())
	if rctx == nil {
		// No route context found
		return "unknown"
	}

	if pattern := rctx.RoutePattern(); pattern != "" {
		// Pattern is already available if it has determined a final route.
		// This happens if the request has already been executed by the router.
		return pattern
	}

	// If the pattern is not available, try to match the route again
	routePath := r.URL.Path
	if r.URL.RawPath != "" {
		routePath = r.URL.RawPath
	}

	tctx := chi.NewRouteContext()
	if !rctx.Routes.Match(tctx, r.Method, routePath) {
		// No match found
		return "unknown"
	}

	// tctx has the updated pattern, since Match mutates it
	return tctx.RoutePattern()
}

// isIgnoredRoute checks if the route should be ignored from monitoring.
func isIgnoredRoute(route string) bool {
	return route == "/health"
}

// monitorLatency measures and records the request latency.
func monitorLatency(method, route string, start time.Time) {
	elapsed := time.Since(start).Seconds()
	monitoring.ObserveLatency(method, route, elapsed)
}
