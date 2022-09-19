package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rainset/shortener/internal/storage"
	"sync"
)

type App struct {
	mutex  sync.RWMutex
	Config Config
	Router *mux.Router
	s      storage.InterfaceStorage
}

type Config struct {
	ServerAddress string
	ServerBaseURL string
}

func New(storage storage.InterfaceStorage) *App {
	return &App{
		s: storage,
	}
}

func (a *App) SetConfigServerAddress(value string) {
	a.Config.ServerAddress = value
}
func (a *App) SetConfigBaseURL(value string) {
	a.Config.ServerBaseURL = value
}
func (a *App) GenerateShortenURL(shortenCode string) string {
	return fmt.Sprintf("%s/%s", a.Config.ServerBaseURL, shortenCode)
}
