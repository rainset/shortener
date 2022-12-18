package main

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-contrib/pprof"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	//gin.SetMode(gin.ReleaseMode)

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
		cnfFileData.BaseURL = *baseURL
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

	srv := &http.Server{
		Addr:    cnfFileData.ServerAddress,
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {

		if cnfFileData.EnableHTTPS == true {
			err = r.RunTLS(cnfFileData.ServerAddress, "cert/cert.pem", "cert/private.key")
		} else {
			err = r.Run(cnfFileData.ServerAddress)
		}
		if err != nil {
			log.Println(cnfFileData)
			panic(err)
		}
		log.Println("Listening on ", cnfFileData.ServerAddress)
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
