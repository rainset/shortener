package main

import (
	"compress/gzip"
	"flag"
	"github.com/gorilla/handlers"
	"github.com/rainset/shortener/internal/app"
	"github.com/rainset/shortener/internal/storage"
	"github.com/rainset/shortener/internal/storage/file"
	"github.com/rainset/shortener/internal/storage/memory"
	"github.com/rainset/shortener/internal/storage/postgres"
	"log"
	"net/http"
	"os"
)

var (
	serverAddress   *string
	baseURL         *string
	fileStoragePath *string
	databaseDsn     *string
)

func init() {
	serverAddress = flag.String("a", os.Getenv("SERVER_ADDRESS"), "string server name, ex:[localhost:8080]")
	baseURL = flag.String("b", os.Getenv("BASE_URL"), "string base url, ex:[http://localhost]")
	fileStoragePath = flag.String("f", os.Getenv("FILE_STORAGE_PATH"), "string file storage path, ex:[/file_storage.log]")
	databaseDsn = flag.String("d", os.Getenv("DATABASE_DSN"), "string db connection, ex:[postgres://root:12345@localhost:5432/shorten]")
}

func main() {

	flag.Parse()

	if *serverAddress == "" {
		*serverAddress = "localhost:8080"
	}
	if *baseURL == "" {
		*baseURL = "http://localhost:8080"
	}

	var s storage.InterfaceStorage

	switch {
	case *databaseDsn != "":
		s = postgres.Init(*databaseDsn)
	case *fileStoragePath != "":
		s = file.Init(*fileStoragePath)
	default:
		s = memory.Init()
	}

	application := app.New(s)
	application.SetConfigServerAddress(*serverAddress)
	application.SetConfigBaseURL(*baseURL)

	r := application.NewRouter()
	http.Handle("/", r)

	log.Printf("Listening %s ...", *serverAddress)
	log.Fatal(http.ListenAndServe(*serverAddress, handlers.CompressHandlerLevel(r, gzip.BestSpeed)))

}
