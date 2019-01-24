package handlers_test

import (
	"encoding/json"
	"errors"
	"goapi/handlers"
	"goapi/models"
	"net/http"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestLoginValidation(t *testing.T) {
	tests := []struct {
		json string
	}{
		{json: "hello"},
		{json: `{"username":"johndoe"}`},
		{json: `{"password":"johndoe"}`},
		{json: `{"username":"","password":""}`},
		{json: `{"username":null,"password":null}`},
	}

	mockUserDataStore := &mockUserDataStore{}

	handler := handlers.LoginHandlerFactory(mockUserDataStore, "signingkey")

	for _, tt := range tests {
		if rr := testHttpRequest(tt.json, handler); rr.Code != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
		}
	}
}

func TestLoginInvalidUsername(t *testing.T) {
	mockUserDataStore := &mockUserDataStore{}
	mockUserDataStore.err = errors.New("User Not Found")

	handler := handlers.LoginHandlerFactory(mockUserDataStore, "signingkey")

	if rr := testHttpRequest(`{"username":"johndoe","password":"pass1234"}`, handler); rr.Code != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusUnauthorized)
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	mockUserDataStore := &mockUserDataStore{}
	mockUserDataStore.user = models.User{
		ID:       100,
		Username: "johndoe",
		Password: "abcdefghij",
		Email:    "johndoe@example.com",
		Active:   true,
		Created:  "2019-01-01 00:00:00",
		Updated:  "2019-01-01 00:00:00",
	}

	handler := handlers.LoginHandlerFactory(mockUserDataStore, "signingkey")

	if rr := testHttpRequest(`{"username":"johndoe","password":"pass1234"}`, handler); rr.Code != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusUnauthorized)
	}
}

func TestLoginValidCredentials(t *testing.T) {
	encPassword, _ := bcrypt.GenerateFromPassword([]byte("pass1234"), 3)

	mockUserDataStore := &mockUserDataStore{}
	mockUserDataStore.user = models.User{
		ID:       100,
		Username: "johndoe",
		Password: string(encPassword),
		Email:    "johndoe@example.com",
		Active:   true,
		Created:  "2019-01-01 00:00:00",
		Updated:  "2019-01-01 00:00:00",
	}

	handler := handlers.LoginHandlerFactory(mockUserDataStore, "signingkey")

	rr := testHttpRequest(`{"username":"johndoe","password":"pass1234"}`, handler)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	var resData map[string]string

	err := json.Unmarshal(rr.Body.Bytes(), &resData)
	if err != nil {
		t.Error(err)
	}

	if len(resData["jwt"]) == 0 {
		t.Error("jwt not set in successful login response")
	}
}
