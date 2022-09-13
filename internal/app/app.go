package app

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/gorilla/mux"
	"github.com/rainset/shortener/internal/app/helper"
	"github.com/rainset/shortener/internal/app/storage/file"
	"github.com/rainset/shortener/internal/app/storage/postgres"
	"log"
	"net/url"
	"sync"
)

type Config struct {
	ServerAddress         string `env:"SERVER_ADDRESS"`
	ServerBaseURL         string `env:"BASE_URL"`
	ServerFileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN           string `env:"DATABASE_DSN"`
	AppKey                string
}

type App struct {
	mutex           sync.RWMutex
	Config          Config
	Router          *mux.Router
	urls            map[string]string
	userHistoryURLs map[string][]string
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

	//if cfg.DatabaseDSN == "" {
	//	cfg.DatabaseDSN = "postgres://root:12345@localhost:5432/shorten"
	//}

	////cfg.AppKey = "49a8aca82c132d8d1f430e32be1e6ff3"
	//cfg.AppKey = "1234567890123456789012345678901234567890"

	//if cfg.ServerStoragePath == "" {
	//	cfg.ServerStoragePath = "storage.log"
	//}

	//db , err := postgres.InitDB(cfg.DatabaseDSN)
	//if err != nil {
	//	fmt.Println(err)
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

	userHistoryURLs := make(map[string][]string)

	return &App{
		Config:          cfg,
		urls:            urls,
		userHistoryURLs: userHistoryURLs,
	}
}

func (a *App) AddURL(value string) (hash string, err error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	urlValue, err := url.ParseRequestURI(value)
	if err != nil {
		return "", err
	}

	hash = helper.GenerateToken(8)
	a.urls[hash] = urlValue.String()

	if a.Config.DatabaseDSN != "" {
		err = postgres.AddURL(hash, urlValue.String())
	}

	if a.Config.ServerFileStoragePath != "" {
		producer, errF := file.NewProducer(a.Config.ServerFileStoragePath)
		if errF != nil {
			return hash, errF
		}
		defer producer.Close()

		requestData := &file.DataURL{Hash: hash, LongURL: value}
		if fileErr := producer.WriteURL(requestData); fileErr != nil {
			return hash, fileErr
		}
	}
	return hash, err
}

func (a *App) AddUserHistoryURL(userID, hash string) (err error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.userHistoryURLs[userID] = append(a.userHistoryURLs[userID], hash)
	return nil
}

func (a *App) GetURL(urlID string) string {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	if a.Config.DatabaseDSN != "" {
		res, err := postgres.GetURL(urlID)
		fmt.Println(res)
		if err == nil && res.Original != "" {
			return res.Original
		}
	}

	return a.urls[urlID]
}

func (a *App) GenerateShortenURL(shortenCode string) string {
	return fmt.Sprintf("%s/%s", a.Config.ServerBaseURL, shortenCode)
}
