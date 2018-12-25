package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"goapi/handlers"
	"goapi/middleware"
	"goapi/models"

	"github.com/gorilla/mux"
)

func main() {
	dsn, noDsn := os.LookupEnv("GOAPI_DSN")
	if noDsn == false {
		log.Fatal("Environment variable GOAPI_DSN must be set.")
	}

	db, err := models.NewDB(dsn)
	if err != nil {
		log.Fatal("Unable to connect with the given DSN.")
	}

	env := handlers.Env{DB: db}

	r := mux.NewRouter()

	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.ContentTypeMiddleware())

	r.HandleFunc("/login", env.LoginHandler).Methods("POST", "OPTIONS")

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      r,
		Addr:         "localhost:8888",
	}

	log.Fatal(srv.ListenAndServe())
}
