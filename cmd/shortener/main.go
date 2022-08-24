package main

import (
	"github.com/gorilla/mux"
	"github.com/rainset/shortener/internal/app"
	"log"
	"net/http"
)

const addr = "localhost:8080"

func main() {

	application := app.New()

	r := mux.NewRouter()
	r.HandleFunc("/", application.SaveURLHandler).Methods("POST")
	r.HandleFunc("/{id:[0-9a-z]+}", application.GetURLHandler).Methods("GET")
	r.HandleFunc("/api/shorten", application.SaveURLJsonHandler).Methods("POST")
	http.Handle("/", r)

	log.Print("Listening...")
	log.Fatal(http.ListenAndServe(addr, r))
}
