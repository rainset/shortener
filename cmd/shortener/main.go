package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"os"

	"github.com/rainset/shortener/internal/app"
	"github.com/rainset/shortener/internal/storage"
	"github.com/rainset/shortener/internal/storage/file"
	"github.com/rainset/shortener/internal/storage/memory"
	"github.com/rainset/shortener/internal/storage/postgres"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"

	serverAddress   *string
	baseURL         *string
	fileStoragePath *string
	databaseDsn     *string
	enableHTTPS     *string
)

func init() {
	// регистрация структуры для сессии
	gob.Register(app.Session{})

	serverAddress = flag.String("a", os.Getenv("SERVER_ADDRESS"), "string server name, ex:[localhost:8080]")
	baseURL = flag.String("b", os.Getenv("BASE_URL"), "string base url, ex:[http://localhost]")
	fileStoragePath = flag.String("f", os.Getenv("FILE_STORAGE_PATH"), "string file storage path, ex:[/file_storage.txt]")
	databaseDsn = flag.String("d", os.Getenv("DATABASE_DSN"), "string db connection, ex:[postgres://root:12345@localhost:5432/shorten]")
	enableHTTPS = flag.String("s", os.Getenv("ENABLE_HTTPS"), "enable https on app")
}

func main() {

	var err error
	var s storage.InterfaceStorage

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
	flag.Parse()

	if *serverAddress == "" {
		*serverAddress = "localhost:8080"
	}
	if *baseURL == "" {
		*baseURL = "http://localhost:8080"
	}

	switch {
	case *databaseDsn != "":
		s = postgres.New(*databaseDsn)
	case *fileStoragePath != "":
		s = file.New(*fileStoragePath)
	default:
		s = memory.New()
	}

	conf := app.Config{
		ServerAddress:  *serverAddress,
		ServerBaseURL:  *baseURL,
		CookieHashKey:  "49a8aca82c132d8d1f430e32be1e6ff3",
		CookieBlockKey: "49a8aca82c132d8d1f430e32be1e6ff2",
	}
	application := app.New(s, conf)
	r := application.NewRouter()

	//pprof.Register(r)

	if *enableHTTPS != "" {
		err = r.RunTLS(conf.ServerAddress, "cert/cert.pem", "cert/private.key")
	} else {
		err = r.Run(conf.ServerAddress)
	}

	if err != nil {
		panic(err)
	}

}
