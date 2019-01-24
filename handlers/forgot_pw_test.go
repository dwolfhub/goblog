package handlers_test

import (
	"errors"
	"goapi/handlers"
	"goapi/models"
	"net/http"
	"testing"
)

func TestForgotPasswordValidation(t *testing.T) {
	tests := []struct {
		json string
	}{
		{json: "hello"},
		{json: `{"hello":"hi"}`},
		{json: `{"hello":"me@example.come"}`},
		{json: `{"email":"hi"}`},
	}

	mockUDS := getMockUserDataStore(errors.New("user not found"), nil)
	mockES := getMockEmailSender(nil, func() {
		t.Error("Email should not have been sent")
	})

	handler := handlers.NewForgotPwHandler(mockUDS, mockES)

	for _, tt := range tests {
		if rr := testHttpRequest(tt.json, handler); rr.Code != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
		}
	}
}

func TestForgotPasswordEmailNotFound(t *testing.T) {
	mockUserDataStore := &mockUserDataStore{}
	mockUserDataStore.err = errors.New("user not found")
	mockEmailSender := &mockEmailSender{}
	mockEmailSender.f = func() {
		t.Error("Email should not have been sent")
	}

	handler := handlers.NewForgotPwHandler(mockUserDataStore, mockEmailSender)

	if rr := testHttpRequest(`{"email":"notfound@example.com"}`, handler); rr.Code != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusCreated)
	}
}

func TestForgotPasswordEmailFound(t *testing.T) {
	mockUserDataStore := &mockUserDataStore{}
	mockUserDataStore.user = models.User{
		ID: 1,
	}
	mockEmailSender := &mockEmailSender{}
	mockEmailSender.f = func() {}

	handler := handlers.NewForgotPwHandler(mockUserDataStore, mockEmailSender)

	if rr := testHttpRequest(`{"email":"notfound@example.com"}`, handler); rr.Code != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusCreated)
	}
}
