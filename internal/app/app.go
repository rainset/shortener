package app

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/mux"
	"github.com/rainset/shortener/internal/app/storage/file"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

	urlValue, err := url.ParseRequestURI(value)
	if err != nil {
		hash = ""
		return
	}

	binHash := md5.Sum([]byte(urlValue.String()))
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

func (a *App) NewRouter() chi.Router {
	r := chi.NewRouter()
	//r.Use(middleware.Compress(5))
	r.Get("/{id:[0-9a-z]+}", a.GetURLHandler)
	r.Post("/api/shorten", a.SaveURLJSONHandler)
	r.Post("/", a.SaveURLHandler)
	r.Get("/{id:[0-9a-z]+}", a.GetURLHandler)
	return r
}

func (a *App) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "id")
	urlValue := a.GetURL(urlID)

	if urlValue == "" {
		http.Error(w, "Bad Url", 400)
		return
	}

	http.Redirect(w, r, urlValue, http.StatusTemporaryRedirect)
}

func (a *App) SaveURLHandler(w http.ResponseWriter, r *http.Request) {
	var bodyBytes []byte
	var err error

	if r.Body != nil {
		bodyBytes, err = ioutil.ReadAll(r.Body)
		if err != nil || len(bodyBytes) == 0 {
			http.Error(w, "Body reading error", 400)
			return
		}
		defer r.Body.Close()
	}
	code, err := a.AddURL(string(bodyBytes))

	if err != nil {
		http.Error(w, fmt.Sprintf("incorrect url format, code: %s body: %s", code, string(bodyBytes)), 400)
		return
	}

	shortenURL := a.GenerateShortenURL(code)

	w.WriteHeader(http.StatusCreated)

	_, writeError := w.Write([]byte(shortenURL))
	if writeError != nil {
		http.Error(w, "response body error", 400)
		return
	}

}
func (a *App) SaveURLJSONHandler(w http.ResponseWriter, r *http.Request) {

	type ShortenRequest struct {
		URL string `json:"url"`
	}
	type ShortenResponse struct {
		Result string `json:"result"`
	}

	var bodyBytes []byte
	var err error

	if r.Body != nil {
		bodyBytes, err = ioutil.ReadAll(r.Body)
		if err != nil || len(bodyBytes) == 0 {
			a.ShowJSONError(w, 400, "Only Json format required in request body")
			return
		}

		defer r.Body.Close()
	}

	value := ShortenRequest{}
	if err := json.Unmarshal(bodyBytes, &value); err != nil {
		a.ShowJSONError(w, 400, "Only Json format required in request body")
		return
	}

	code, err := a.AddURL(value.URL)

	fmt.Println("code:", code)
	fmt.Println("err:", err)

	if err != nil {
		http.Error(w, "incorrect url format", 400)
		return
	}

	shortenURL := a.GenerateShortenURL(code)
	shortenData := ShortenResponse{Result: shortenURL}
	shortenJSON, err := json.Marshal(shortenData)
	if err != nil {
		http.Error(w, "json response error", 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, writeError := w.Write(shortenJSON)
	if writeError != nil {
		http.Error(w, "response body error", 400)
		return
	}

}

func (a *App) GenerateShortenURL(shortenCode string) string {
	return fmt.Sprintf("%s/%s", a.Config.ServerBaseURL, shortenCode)
}

func (a *App) ShowJSONError(w http.ResponseWriter, code int, message string) {

	type ErrorResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	data, err := json.Marshal(ErrorResponse{Code: code, Message: message})
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, writeError := w.Write(data)
	if writeError != nil {
		panic(writeError)
	}

}
