package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
)

// TestHandler is a helper method to quickly test a request handler
func TestHandler(body string, handler func(w http.ResponseWriter, r *http.Request)) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	method := "GET"
	if len(body) > 0 {
		method = "POST"
	}

	req, _ := http.NewRequest(method, "/", bytes.NewBufferString(body))

	r := mux.NewRouter()
	r.HandleFunc("/", handler)

	r.ServeHTTP(rr, req)

	return rr
}
