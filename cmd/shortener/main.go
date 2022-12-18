package main

import (
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-contrib/pprof"
	"log"
	"os"

	"github.com/rainset/shortener/internal/app"
	"github.com/rainset/shortener/internal/storage"
	"github.com/rainset/shortener/internal/storage/file"
	"github.com/rainset/shortener/internal/storage/memory"
	"github.com/rainset/shortener/internal/storage/postgres"
)

//go:generate go run ../certificate/certificate.go

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"

	serverAddress   *string
	baseURL         *string
	fileStoragePath *string
	databaseDsn     *string
	enableHTTPS     *string
	configFile      *string
)

type ConfigFileData struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	Database        string `json:"database"`
	EnableHTTPS     bool   `json:"enable_https"`
}

func init() {
	// регистрация структуры для сессии
	gob.Register(app.Session{})

	configFile = flag.String("c", os.Getenv("CONFIG"), "path to config.json file")
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

	var cnfFileData ConfigFileData

	if *configFile != "" {
		var errCnf error
		cnfFile, errCnf := os.ReadFile(*configFile)
		if errCnf != nil {
			log.Println("Error when opening file: ", errCnf)
		}
		errCnf = json.Unmarshal(cnfFile, &cnfFileData)
		if errCnf != nil {
			log.Println("Error during Unmarshal(): ", errCnf)
		}
	}

	if *serverAddress != "" {
		cnfFileData.ServerAddress = *serverAddress
	}
	if *baseURL != "" {
		cnfFileData.BaseURL = *serverAddress
	}
	if *baseURL != "" {
		cnfFileData.BaseURL = *serverAddress
	}
	if *fileStoragePath != "" {
		cnfFileData.FileStoragePath = *fileStoragePath
	}
	if *databaseDsn != "" {
		cnfFileData.Database = *databaseDsn
	}
	if *enableHTTPS != "" {
		cnfFileData.EnableHTTPS = true
	}

	switch {
	case cnfFileData.Database != "":
		s = postgres.New(cnfFileData.Database)
	case cnfFileData.FileStoragePath != "":
		s = file.New(cnfFileData.FileStoragePath)
	default:
		s = memory.New()
	}

	conf := app.Config{
		ServerAddress:  cnfFileData.ServerAddress,
		ServerBaseURL:  cnfFileData.BaseURL,
		CookieHashKey:  "49a8aca82c132d8d1f430e32be1e6ff3",
		CookieBlockKey: "49a8aca82c132d8d1f430e32be1e6ff2",
	}
	application := app.New(s, conf)
	r := application.NewRouter()

	pprof.Register(r)

	if cnfFileData.EnableHTTPS == true {
		err = r.RunTLS(cnfFileData.ServerAddress, "cert/cert.pem", "cert/private.key")
	} else {
		err = r.Run(cnfFileData.ServerAddress)
	}

	if err != nil {
		panic(err)
	}

}
