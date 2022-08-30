package main

import (
	"github.com/rainset/shortener/internal/app"
	"log"
	"net/http"
)

func main() {
	application := app.New()
	application.InitFlags()
	r := application.NewRouter()
	log.Printf("Listening %s ...", application.Config.ServerAddress)
	log.Fatal(http.ListenAndServe(application.Config.ServerAddress, r))
}
