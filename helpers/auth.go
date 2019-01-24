package helpers

import (
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// AppClaims defines the custom claims for our JWTs
type AppClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateSignedString generates a JWT string
func GenerateSignedString(username string, signingKey string) string {
	claims := AppClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			Issuer:    "website",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, err := token.SignedString([]byte(signingKey))

	if err != nil {
		log.Fatal("Unable to generate token using signing key")
	}

	return ss
}
