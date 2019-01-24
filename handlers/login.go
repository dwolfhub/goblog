package handlers

import (
	"encoding/json"
	"goapi/helpers"
	"goapi/models"
	"net/http"
	"net/url"

	"golang.org/x/crypto/bcrypt"

	"github.com/thedevsaddam/govalidator"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type jwtResponse struct {
	Jwt string `json:"jwt"`
}

// LoginHandlerFactory creates a func that handles login requests
func LoginHandlerFactory(userStore models.UserDataReader, signingKey string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		creds, errs := validateLoginRequest(w, req)

		if len(errs) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"errors": errs})
			return
		}

		user, err := userStore.GetUserByUsername(creds.Username)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ss := helpers.GenerateSignedString(user.Username, signingKey)

		json.NewEncoder(w).Encode(jwtResponse{
			Jwt: ss,
		})
	}
}

func validateLoginRequest(w http.ResponseWriter, req *http.Request) (credentials, url.Values) {
	var credentials credentials

	rules := govalidator.MapData{
		"username": []string{"required"},
		"password": []string{"required"},
	}

	messages := govalidator.MapData{
		"username": []string{"required:This field is required."},
		"password": []string{"required:This field is required."},
	}

	opts := govalidator.Options{
		Request:         req,
		Data:            &credentials,
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()

	return credentials, e
}
