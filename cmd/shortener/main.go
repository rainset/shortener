package main

import (
	"github.com/klauspost/compress/gzhttp"
	"github.com/rainset/shortener/internal/app"
	"log"
	"net/http"
)

func main() {
	application := app.New()
	application.InitFlags()
	r := application.NewRouter()
	http.Handle("/", gzhttp.GzipHandler(r))
	log.Printf("Listening %s ...", application.Config.ServerAddress)
	log.Fatal(http.ListenAndServe(application.Config.ServerAddress, r))
}
