package handlers_test

import (
	"encoding/json"
	"goapi/handlers"
	"goapi/helpers"
	"goapi/middleware"
	"goapi/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetMeNoToken(t *testing.T) {
	mockUserDataStore := &mockUserDataStore{}

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/user/me", handlers.GetMeHandlerFactory(mockUserDataStore))

	req, _ := http.NewRequest("GET", "/user/me", nil)

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestGetMeReturnsUserData(t *testing.T) {
	const signingKey string = "signingKey"

	user := models.User{
		ID:       100,
		Username: "johndoe",
		Password: "abcdefghij",
		Email:    "johndoe@example.com",
		Active:   true,
		Created:  "2019-01-01 00:00:00",
		Updated:  "2019-01-01 00:00:00",
	}

	mockUserDataStore := &mockUserDataStore{
		user: user,
	}

	ss := helpers.GenerateSignedString(user.Username, signingKey)

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.Use(middleware.AuthRequiredMiddleware(signingKey))
	r.HandleFunc("/user/me", handlers.GetMeHandlerFactory(mockUserDataStore))

	req, _ := http.NewRequest("GET", "/user/me", nil)
	req.Header.Set("Authorization", strings.Join([]string{"Bearer: ", ss}, ""))

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resData map[string]interface{}

	err := json.Unmarshal(rr.Body.Bytes(), &resData)
	if err != nil {
		t.Error(err)
	}

	if id, isInt := resData["id"].(int); isInt && id != 100 {
		t.Error("id not set successfully in login response")
	}
	if username, isString := resData["username"].(string); isString && username != "johndoe" {
		t.Error("username not set successfully in login response")
	}
	if email, isString := resData["email"].(string); isString && email != "johndoe@example.com" {
		t.Error("email not set successfully in login response")
	}
}
