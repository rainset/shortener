package main

import (
	"github.com/gorilla/handlers"
	"github.com/rainset/shortener/internal/app"
	"log"
	"net/http"
)

func main() {

	application := app.New()
	application.InitFlags()

	r := application.NewRouter()
	http.Handle("/", r)

	log.Printf("Listening %s ...", application.Config.ServerAddress)
	log.Fatal(http.ListenAndServe(application.Config.ServerAddress, handlers.CompressHandler(r)))
}
