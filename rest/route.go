package rest

import (
	"net/http"

	"github.com/dexlabsio/garlic/errors"
	"github.com/dexlabsio/garlic/request"
)

type RouteOptions int

const (
	// RouteOptionAllowPublicAccess tells the system to allow public access
	// to this route without authentication
	RouteOptionAllowPublicAccess RouteOptions = iota

	// RouteOptionAllowOnlyAuthenticated tells the system to block the specific
	// route for external access. Only authenticated users can access this. This
	// is the default behavior
	RouteOtionAllowOnlyAuthenticated

	// RouteOptionAllowOnlySuperuser tells the system to block all the access to
	// this route except if you're accessing this authenticated as a superuser
	RouteOptionAllowOnlySuperuser
)

type RouteOptionsMap map[RouteOptions]struct{}

type Route struct {
	URL          string
	Method       string
	Handler      http.HandlerFunc
	RouteOptions RouteOptionsMap
}

type App interface {
	Routes() []*Route
}

func (o RouteOptionsMap) IsEmpty() bool {
	return len(o) == 0
}

func ToHandler(f func(http.ResponseWriter, *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			l := request.GetLogger(r)

			if errors.IsKind(err, errors.KindUserError) {
				l.Warn("[USER ERROR]", errors.Zap(err))
			} else {
				l.Error("[SYSTEM ERROR]", errors.Zap(err))
			}

			WriteError(err).Must(w)
		}
	}
}

func route(method string, url string, f func(http.ResponseWriter, *http.Request) error, opts ...RouteOptions) *Route {
	handler := ToHandler(f)

	optsMap := make(RouteOptionsMap, len(opts))
	for _, opt := range opts {
		optsMap[opt] = struct{}{}
	}

	return &Route{
		URL:          url,
		Method:       method,
		Handler:      handler,
		RouteOptions: optsMap,
	}
}

func Get(url string, f func(http.ResponseWriter, *http.Request) error, opts ...RouteOptions) *Route {
	return route(http.MethodGet, url, f, opts...)
}

func Post(url string, f func(http.ResponseWriter, *http.Request) error, opts ...RouteOptions) *Route {
	return route(http.MethodPost, url, f, opts...)
}

func Put(url string, f func(http.ResponseWriter, *http.Request) error, opts ...RouteOptions) *Route {
	return route(http.MethodPut, url, f, opts...)
}

func Patch(url string, f func(http.ResponseWriter, *http.Request) error, opts ...RouteOptions) *Route {
	return route(http.MethodPatch, url, f, opts...)
}

func Delete(url string, f func(http.ResponseWriter, *http.Request) error, opts ...RouteOptions) *Route {
	return route(http.MethodDelete, url, f, opts...)
}
