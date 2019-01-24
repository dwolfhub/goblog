package handlers

import (
	"encoding/json"
	"goapi/helpers"
	"goapi/models"
	"net/http"
	"net/url"

	"github.com/thedevsaddam/govalidator"
)

type forgotPwRequest struct {
	Email string `json:"email"`
}

// NewForgotPwHandler returns a method that handles forgot password requests
func NewForgotPwHandler(userStore models.UserDataReader, emailSender helpers.EmailSender) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestData, errs := validateForgotPwRequest(w, r)

		if len(errs) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"errors": errs})
			return
		}

		user, err := userStore.GetUserByEmail(requestData.Email)

		if err == nil {
			emailSender.EmailSend([]string{user.Email}, "hello!", "hello!!!")
		}

		w.WriteHeader(201)
	}
}

func validateForgotPwRequest(w http.ResponseWriter, req *http.Request) (forgotPwRequest, url.Values) {
	var forgotPwRequest forgotPwRequest

	rules := govalidator.MapData{
		"email": []string{"required", "email"},
	}

	messages := govalidator.MapData{
		"email": []string{"required:This field is required.", "email:Please enter a valid email."},
	}

	opts := govalidator.Options{
		Request:         req,
		Data:            &forgotPwRequest,
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()

	return forgotPwRequest, e
}
