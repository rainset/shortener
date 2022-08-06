// Package app является основным приложение, в котором реализована инициализация бд и http хендлеров
package app

import (
	"context"
	"fmt"
	"github.com/gin-contrib/pprof"
	pb "github.com/rainset/shortener/internal/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rainset/shortener/internal/queue"
	"github.com/rainset/shortener/internal/storage"
)

type App struct {
	dq     *queue.DeleterQueue
	s      storage.InterfaceStorage
	Config Config
}

type Config struct {
	ServerAddress  string
	ServerBaseURL  string
	CookieHashKey  string
	CookieBlockKey string
	TrustedSubnet  string
	EnableHTTPS    bool
	GRPCMode       bool
}

// Session структура для хранения сессии пользователя
type Session struct {
	UserID string
}

// New инициализация основного объекта приложения
func New(storage storage.InterfaceStorage, c Config) *App {

	dq := queue.NewDeleterQueue(storage)
	go dq.Init()

	return &App{
		s:      storage,
		Config: c,
		dq:     dq,
	}

}

func (a *App) StartHTTPServer() {

	shs := &ShortenerHTTPServer{
		a: a,
	}

	r := shs.NewRouter()

	pprof.Register(r)

	srv := &http.Server{
		Addr:    a.Config.ServerAddress,
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		var err error
		if a.Config.EnableHTTPS {
			err = r.RunTLS(a.Config.ServerAddress, "cert/cert.pem", "cert/private.key")
		} else {
			err = r.Run(a.Config.ServerAddress)
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Listening on ", a.Config.ServerAddress)
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 3)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
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

func (a *App) StartGRPCServer() {
	// определяем порт для сервера
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()
	// регистрируем сервис
	pb.RegisterShortenerServer(s, &ShortenerGRPCServer{a: a})

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
