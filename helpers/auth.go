package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func init() {
	assertAvailablePRNG()
}

func assertAvailablePRNG() {
	// Assert that a cryptographically secure PRNG is available.
	// Panic otherwise.
	buf := make([]byte, 1)

	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: Read() failed with %#v", err))
	}
}

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

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomToken will generate a URL-safe random string for use as a password reset token or similar
func GenerateRandomToken() (string, error) {
	bytes, err := generateRandomBytes(32)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), err
}
