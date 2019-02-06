package auth_test

import (
	"encoding/json"
	"errors"
	"goapi/handlers/auth"
	"goapi/models"
	"net/http"
	"testing"

	"goapi/handlers"

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

	mockUserDataStore := &models.MockUserDataStore{}

	handler := auth.LoginHandlerFactory(mockUserDataStore, "signingkey")

	for _, tt := range tests {
		if rr := handlers.TestHandler(tt.json, handler); rr.Code != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
		}
	}
}

func TestLoginInvalidUsername(t *testing.T) {
	mockUserDataStore := &models.MockUserDataStore{}
	mockUserDataStore.Err = errors.New("User Not Found")

	handler := auth.LoginHandlerFactory(mockUserDataStore, "signingkey")

	if rr := handlers.TestHandler(`{"username":"johndoe","password":"pass1234"}`, handler); rr.Code != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusUnauthorized)
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	mockUserDataStore := &models.MockUserDataStore{}
	mockUserDataStore.User = models.User{
		ID:       100,
		Username: "johndoe",
		Password: "abcdefghij",
		Email:    "johndoe@example.com",
		Active:   true,
		Created:  "2019-01-01 00:00:00",
		Updated:  "2019-01-01 00:00:00",
	}

	handler := auth.LoginHandlerFactory(mockUserDataStore, "signingkey")

	if rr := handlers.TestHandler(`{"username":"johndoe","password":"pass1234"}`, handler); rr.Code != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusUnauthorized)
	}
}

func TestLoginValidCredentials(t *testing.T) {
	encPassword, _ := bcrypt.GenerateFromPassword([]byte("pass1234"), 3)

	mockUserDataStore := &models.MockUserDataStore{}
	mockUserDataStore.User = models.User{
		ID:       100,
		Username: "johndoe",
		Password: string(encPassword),
		Email:    "johndoe@example.com",
		Active:   true,
		Created:  "2019-01-01 00:00:00",
		Updated:  "2019-01-01 00:00:00",
	}

	handler := auth.LoginHandlerFactory(mockUserDataStore, "signingkey")

	rr := handlers.TestHandler(`{"username":"johndoe","password":"pass1234"}`, handler)

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
