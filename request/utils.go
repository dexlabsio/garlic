package request

import (
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
)

func RouteContainsPattern(r *http.Request, pattern string) bool {
	// Get the route pattern for the current request
	routePattern := chi.RouteContext(r.Context()).RoutePattern()

	// Define the regex pattern to match {organization_id} in the route
	re := regexp.MustCompile(pattern)
	return re.MatchString(routePattern)
}
