package middleware

import (
	"goapi/handlers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"
)

func TestBadAuthorizationHeaderReturnsUnauthorized(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.Use(AuthRequiredMiddleware("securitykey"))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	badTokens := []string{
		"",
		"Bearer: ",
		"BEARER: ASADF",
		"Bearer :ASDFSDF",
	}

	for _, badToken := range badTokens {
		req.Header.Add("Authorization", badToken)

		r.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	}
}

func TestGoodAuthorizationHeaderCallsHandler(t *testing.T) {
	securityKey := "securitykey"

	req, _ := http.NewRequest("GET", "/", nil)

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.Use(AuthRequiredMiddleware(securityKey))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	claims := handlers.AppClaims{
		Username: "johndoe123",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15000).Unix(),
			Issuer:    "website",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, _ := token.SignedString([]byte(securityKey))

	req.Header.Add("Authorization", strings.Join([]string{"Bearer: ", ss}, ""))

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}
