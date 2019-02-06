package main

import (
	"crypto/tls"
	"goapi/handlers"
	"goapi/handlers/auth"
	"goapi/helpers"
	"log"
	"net/http"
	"time"

	"goapi/middleware"
	"goapi/models"

	"github.com/gorilla/mux"
)

var (
	db          *models.DB
	signingKey  string
	emailSender helpers.EmailSender
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	signingKey = helpers.GetEnvVar("GOAPI_SIGNING_KEY")
	dsn := helpers.GetEnvVar("GOAPI_DSN")

	var err error
	db, err = models.NewDB(dsn)
	if err != nil {
		log.Fatal("Unable to connect with the given DSN.")
	}

	emailSender = helpers.NewEmailSender(helpers.EmailServerConfigurer{
		EmailHost:     helpers.GetEnvVar("GOAPI_EMAIL_HOST"),
		EmailPort:     helpers.GetEnvVar("GOAPI_EMAIL_PORT"),
		EmailFrom:     helpers.GetEnvVar("GOAPI_EMAIL_FROM"),
		EmailLogin:    helpers.GetEnvVar("GOAPI_EMAIL_LOGIN"),
		EmailPassword: helpers.GetEnvVar("GOAPI_EMAIL_PASSWORD"),
	})
}

func main() {
	mainRouter := mux.NewRouter()
	mainRouter.Use(middleware.CorsMiddleware())
	mainRouter.Use(middleware.ContentTypeMiddleware())
	mainRouter.HandleFunc("/login", auth.LoginHandlerFactory(db, signingKey)).Methods("POST", "OPTIONS")
	mainRouter.HandleFunc("/forgot", auth.NewForgotPwHandler(db, emailSender)).Methods("POST", "OPTIONS")
	mainRouter.HandleFunc("/post", handlers.GetPostsHandlerFactory(db)).Methods("GET", "OPTIONS")

	userRouter := mainRouter.PathPrefix("/user").Subrouter()
	userRouter.Use(middleware.AuthRequiredMiddleware(signingKey))
	userRouter.HandleFunc("/me", handlers.GetMeHandlerFactory(db)).Methods("GET", "OPTIONS")

	srv := &http.Server{
		Addr:         ":443",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mainRouter,
		TLSConfig: &tls.Config{
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519,
			},
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		},
	}

	go http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		target := "https://" + req.Host + req.URL.Path
		if len(req.URL.RawQuery) > 0 {
			target += "?" + req.URL.RawQuery
		}
		http.Redirect(w, req, target, http.StatusTemporaryRedirect)
	}))

	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}
