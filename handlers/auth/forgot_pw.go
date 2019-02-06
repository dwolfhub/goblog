package auth

import (
	"encoding/json"
	"fmt"
	"goapi/helpers"
	"goapi/models"
	"log"
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
			token, err := helpers.GenerateRandomToken()
			if err != nil {
				log.Fatalf("Unable to generate random token: %s", err)
			}

			emailSender.EmailSend([]string{user.Email}, "Reset Password Request", fmt.Sprintf("https://%s/reset?t=%s", r.Host, token))
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
