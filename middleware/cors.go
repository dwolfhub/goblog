package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

// CorsMiddleware adds CORS headers
func CorsMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if req.Method == "OPTIONS" {
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}
