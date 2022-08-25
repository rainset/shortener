package main

import (
	"github.com/gorilla/mux"
	"github.com/rainset/shortener/internal/app"
	"log"
	"net/http"
)

func main() {
	application := app.New()

	r := mux.NewRouter()
	r.HandleFunc("/{id:[0-9a-z]+}", application.GetURLHandler).Methods("GET")
	r.HandleFunc("/api/shorten", application.SaveURLJSONHandler).Methods("POST")
	r.HandleFunc("/", application.SaveURLHandler).Methods("POST")
	http.Handle("/", r)

	log.Printf("Listening %s ...", application.Config.ServerAddress)
	log.Fatal(http.ListenAndServe(application.Config.ServerAddress, r))
}
