package rest

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"

	"github.com/dexlabsio/garlic/errors"
	"github.com/dexlabsio/garlic/logging"
)

type Routes []*Route

type Route struct {
	Method  string
	Pattern string
	Fn      func(http.ResponseWriter, *http.Request) error
}

func (route *Route) Handler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := route.Fn(w, r)
		if err != nil {
			ctx := r.Context()
			l := logging.GetLoggerFromContext(ctx)

			if errors.IsKind(err, errors.KindUserError) {
				l.Warn("[USER ERROR]", errors.Zap(err))
			} else {
				l.Error("[SYSTEM ERROR]", errors.Zap(err))
			}

			WriteError(err).Must(w)
		}
	}
}

func Get(url string, f func(http.ResponseWriter, *http.Request) error) *Route {
	return &Route{http.MethodGet, url, f}
}

func Post(url string, f func(http.ResponseWriter, *http.Request) error) *Route {
	return &Route{http.MethodPost, url, f}
}

func Put(url string, f func(http.ResponseWriter, *http.Request) error) *Route {
	return &Route{http.MethodPut, url, f}
}

func Patch(url string, f func(http.ResponseWriter, *http.Request) error) *Route {
	return &Route{http.MethodPatch, url, f}
}

func Delete(url string, f func(http.ResponseWriter, *http.Request) error) *Route {
	return &Route{http.MethodDelete, url, f}
}

type App interface {
	Routes() Routes
}

func RegisterApp(r chi.Router, app App) {
	for _, route := range app.Routes() {
		r.MethodFunc(route.Method, route.Pattern, route.Handler())
	}
}
