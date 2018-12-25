package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

// ContentTypeMiddleware adds standard headers
func ContentTypeMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			next.ServeHTTP(w, req)
		})
	}
}
