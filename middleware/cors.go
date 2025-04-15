package middleware

import (
	"net/http"
	"strings"
)

func Cors(config *Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", strings.Join(config.Cors.AllowedHosts, ", "))
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.Cors.AllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.Cors.AllowedHeaders, ", "))
			w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.Cors.ExposedHeaders, ", "))

			if r.Method == "OPTIONS" {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
