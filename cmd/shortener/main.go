package main

import (
	"github.com/rainset/shortener/internal/app"
	"log"
	"net/http"
)

const addr = "localhost:8080"

func main() {
	application := app.New()
	http.HandleFunc("/", application.RouteHandler)
	log.Print("Listening...")
	log.Fatal(http.ListenAndServe(addr, nil))
}
