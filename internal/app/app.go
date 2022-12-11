// Package app является основным приложение, в котором реализована инициализация бд и http хендлеров
package app

import (
	"fmt"
	"runtime"

	"github.com/gorilla/mux"

	"github.com/rainset/shortener/internal/queue"
	"github.com/rainset/shortener/internal/storage"
)

type App struct {
	Router *mux.Router
	Queue  *queue.DeleteURLQueue
	s      storage.InterfaceStorage
	Config Config
}

type Config struct {
	ServerAddress  string
	ServerBaseURL  string
	CookieHashKey  string
	CookieBlockKey string
}

// Session структура для хранения сессии пользователя
type Session struct {
	UserID string
}

// New инициализация основного объекта приложения
func New(storage storage.InterfaceStorage, c Config) *App {

	newQueue := queue.NewDeleteURLQueue(storage)
	go newQueue.PeriodicURLDelete()

	workers := make([]*queue.DeleteURLWorker, 0, runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		workers = append(workers, queue.NewDeleteURLWorker(i, newQueue, storage))

	}

	for _, w := range workers {
		go w.Loop()
	}

	return &App{
		s:      storage,
		Config: c,
		Queue:  newQueue,
	}

}

// GenerateShortenURL генерирует web ссылку
func (a *App) GenerateShortenURL(shortenCode string) string {
	return fmt.Sprintf("%s/%s", a.Config.ServerBaseURL, shortenCode)
}
