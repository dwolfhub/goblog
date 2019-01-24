package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestCorsHeadersAdded(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.Use(CorsMiddleware())
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	r.ServeHTTP(rr, req)

	if methods := rr.HeaderMap.Get("Access-Control-Allow-Methods"); methods != "GET,POST,PUT,DELETE,OPTIONS" {
		t.Error("access-control-allow-methods not set properly")
	}

	if origin := rr.HeaderMap.Get("Access-Control-Allow-Origin"); origin != "http://localhost:3000" {
		t.Error("access-control-allow-origin not set properly")
	}
}

func TestOptionsResponseEmpty(t *testing.T) {
	req, _ := http.NewRequest("OPTIONS", "/", nil)

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.Use(CorsMiddleware())
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// should not run
		w.WriteHeader(http.StatusCreated)
	})

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Error("options request not intercepted")
	}
}
