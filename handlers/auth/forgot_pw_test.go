package auth_test

import (
	"errors"
	"goapi/handlers"
	"goapi/handlers/auth"
	"goapi/helpers"
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

	mockUDS := models.GetMockUserDataStore(errors.New("user not found"), nil)
	mockES := helpers.GetMockEmailSender(nil, func() {
		t.Error("Email should not have been sent")
	})

	handler := auth.NewForgotPwHandler(mockUDS, mockES)

	for _, tt := range tests {
		if rr := handlers.TestHandler(tt.json, handler); rr.Code != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
		}
	}
}

func TestForgotPasswordEmailNotFound(t *testing.T) {
	mockUserDataStore := &models.MockUserDataStore{}
	mockUserDataStore.Err = errors.New("user not found")
	mockEmailSender := &helpers.MockEmailSender{}
	mockEmailSender.F = func() {
		t.Error("Email should not have been sent")
	}

	handler := auth.NewForgotPwHandler(mockUserDataStore, mockEmailSender)

	if rr := handlers.TestHandler(`{"email":"notfound@example.com"}`, handler); rr.Code != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusCreated)
	}
}

func TestForgotPasswordEmailFound(t *testing.T) {
	mockUserDataStore := &models.MockUserDataStore{}
	mockUserDataStore.User = models.User{
		ID: 1,
	}
	mockEmailSender := &helpers.MockEmailSender{}
	mockEmailSender.F = func() {}

	handler := auth.NewForgotPwHandler(mockUserDataStore, mockEmailSender)

	if rr := handlers.TestHandler(`{"email":"notfound@example.com"}`, handler); rr.Code != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusCreated)
	}
}
