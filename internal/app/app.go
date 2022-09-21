package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rainset/shortener/internal/cookie"
	"github.com/rainset/shortener/internal/storage"
	"sync"
)

type App struct {
	mutex  sync.RWMutex
	Config Config
	Router *mux.Router
	s      storage.InterfaceStorage
	cookie *cookie.SCookie
}

type Config struct {
	ServerAddress  string
	ServerBaseURL  string
	CookieHashKey  string
	CookieBlockKey string
}

func New(storage storage.InterfaceStorage, c Config) *App {
	return &App{
		s:      storage,
		cookie: cookie.New(c.CookieHashKey, c.CookieBlockKey),
		Config: c,
	}
}

func (a *App) GenerateShortenURL(shortenCode string) string {

	fmt.Println(a.Config)

	return fmt.Sprintf("%s/%s", a.Config.ServerBaseURL, shortenCode)
}
