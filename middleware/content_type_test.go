package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestAddsContentType(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.Use(ContentTypeMiddleware())
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	r.ServeHTTP(rr, req)

	if header := rr.HeaderMap.Get("Content-Type"); header != "application/json" {
		t.Error("content-type header not set")
	}
}
