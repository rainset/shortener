package app

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/caarlos0/env/v6"
	"github.com/gorilla/mux"
	"github.com/rainset/shortener/internal/app/storage/file"
	"log"
	"sync"
)

type Config struct {
	ServerAddress         string `env:"SERVER_ADDRESS"`
	ServerBaseURL         string `env:"BASE_URL"`
	ServerFileStoragePath string `env:"FILE_STORAGE_PATH"`
}

type App struct {
	sync.RWMutex
	Config Config
	Router *mux.Router
	urls   map[string]string
}

func New() *App {

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.ServerAddress == "" {
		cfg.ServerAddress = "localhost:8080"
	}
	if cfg.ServerBaseURL == "" {
		cfg.ServerBaseURL = "http://localhost:8080"
	}
	//if cfg.ServerStoragePath == "" {
	//	cfg.ServerStoragePath = "storage.log"
	//}

	urls := make(map[string]string)

	if cfg.ServerFileStoragePath != "" {
		consumer, err := file.NewConsumer(cfg.ServerFileStoragePath)
		if err != nil {
			log.Fatal(err)
		}
		urls, err = consumer.RestoreStorage()
		if err != nil {
			log.Fatal(err)
		}
	}

	return &App{urls: urls, Config: cfg}
}

func (a *App) AddURL(value string) (hash string, err error) {
	a.RLock()
	defer a.RUnlock()

	//urlValue, err := url.ParseRequestURI(value)
	//if err != nil {
	//	hash = ""
	//	return
	//}

	binHash := md5.Sum([]byte(value))
	hash = hex.EncodeToString(binHash[:])
	a.urls[hash] = value

	if a.Config.ServerFileStoragePath != "" {
		producer, errF := file.NewProducer(a.Config.ServerFileStoragePath)
		if errF != nil {
			return
		}
		defer producer.Close()

		requestData := &file.DataURL{Hash: hash, LongURL: value}
		if fileErr := producer.WriteURL(requestData); fileErr != nil {
			return
		}
	}
	return hash, err
}

func (a *App) GetURL(urlID string) string {
	a.RLock()
	defer a.RUnlock()
	return a.urls[urlID]
}
