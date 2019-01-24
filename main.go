package main

import (
	"goapi/handlers"
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
	mainRouter.HandleFunc("/login", handlers.LoginHandlerFactory(db, signingKey)).Methods("POST", "OPTIONS")
	mainRouter.HandleFunc("/forgot", handlers.NewForgotPwHandler(db, emailSender)).Methods("POST")
	mainRouter.HandleFunc("/post", handlers.GetPostsHandlerFactory(db)).Methods("GET", "OPTIONS")

	userRouter := mainRouter.PathPrefix("/user").Subrouter()
	userRouter.Use(middleware.AuthRequiredMiddleware(signingKey))
	userRouter.HandleFunc("/me", handlers.GetMeHandlerFactory(db)).Methods("GET", "OPTIONS")

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mainRouter,
		Addr:         "localhost:8888",
		// todo TLS
	}

	log.Fatal(srv.ListenAndServe())
}
