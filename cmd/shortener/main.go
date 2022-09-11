package main

import (
	"compress/gzip"
	"github.com/gorilla/handlers"
	"github.com/rainset/shortener/internal/app"
	"github.com/rainset/shortener/internal/app/storage/postgres"
	"log"
	"net/http"
)

func main() {

	application := app.New()
	application.InitFlags()

	if application.Config.DatabaseDSN != "" {
		errDB := postgres.InitDB(application.Config.DatabaseDSN)
		if errDB != nil {
			log.Println(errDB)
		}
		errDB = postgres.CreateTables()
		if errDB != nil {
			log.Println(errDB)
		}
	}

	r := application.NewRouter()
	http.Handle("/", r)

	log.Printf("Listening %s ...", application.Config.ServerAddress)
	log.Fatal(http.ListenAndServe(application.Config.ServerAddress, handlers.CompressHandlerLevel(r, gzip.BestSpeed)))

}
