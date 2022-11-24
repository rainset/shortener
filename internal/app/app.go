package app

import (
	"fmt"
	"runtime"

	"github.com/gorilla/mux"

	"github.com/rainset/shortener/internal/queue"
	"github.com/rainset/shortener/internal/storage"
)

type App struct {
	Config Config
	Router *mux.Router
	s      storage.InterfaceStorage
	Queue  *queue.DeleteURLQueue
}

type Config struct {
	ServerAddress  string
	ServerBaseURL  string
	CookieHashKey  string
	CookieBlockKey string
}

type Session struct {
	UserID string
}

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

func (a *App) GenerateShortenURL(shortenCode string) string {
	return fmt.Sprintf("%s/%s", a.Config.ServerBaseURL, shortenCode)
}
