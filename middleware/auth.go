package middleware

import (
	"goapi/helpers"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// AuthRequiredMiddleware requires that the user be authenticated with JWT
func AuthRequiredMiddleware(secretKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			header := req.Header.Get("Authorization")
			splitHeader := strings.Split(header, ": ")

			if len(splitHeader) != 2 || splitHeader[0] != "Bearer" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			tokenString := strings.Trim(splitHeader[1], " ")

			if len(tokenString) == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			token, _ := jwt.ParseWithClaims(tokenString, &helpers.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			})

			if !token.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(*helpers.AppClaims); ok {
				context.Set(req, "username", claims.Username)
				next.ServeHTTP(w, req)
				return
			}

			w.WriteHeader(http.StatusUnauthorized)
			return
		})
	}
}
