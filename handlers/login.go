package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/thedevsaddam/govalidator"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AppClaims defines the custom claims for our JWTs
type AppClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type jwtResponse struct {
	Jwt string `json:"jwt"`
}

// LoginHandler handles login requests
func (env *Env) LoginHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")

	creds, errs := validate(w, req)

	if len(errs) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"errors": errs})
		return
	}

	user, err := env.DB.GetUserByUsername(creds.Username)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ss := generateSignedString(user.Username)

	json.NewEncoder(w).Encode(jwtResponse{
		Jwt: ss,
	})
}

func generateSignedString(username string) string {
	mySigningKey, isSet := os.LookupEnv("GOAPI_SECURITY_KEY")

	if !isSet {
		log.Fatal("GOAPI_SECURITY_KEY is not set")
	}

	claims := AppClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15000).Unix(),
			Issuer:    "website",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, err := token.SignedString([]byte(mySigningKey))

	if err != nil {
		log.Fatal("Unable to generate token using signing key")
	}

	return ss
}

func validate(w http.ResponseWriter, req *http.Request) (credentials, url.Values) {
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
