package app

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"sync"
)

type App struct {
	sync.RWMutex
	urls map[string]string
}

type dataJSON struct {
	URL string `json:"url"`
}

type ErrorResponse struct {
	Code    int
	Message string
}

func New() *App {
	return &App{urls: make(map[string]string)}
}

func (a *App) AddURL(value string) string {
	a.RLock()
	defer a.RUnlock()
	binHash := md5.Sum([]byte(value))
	hash := hex.EncodeToString(binHash[:])
	a.urls[hash] = value
	return hash
}

func (a *App) GetURL(urlID string) string {
	a.RLock()
	defer a.RUnlock()
	return a.urls[urlID]
}

func (a *App) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urlID := vars["id"]
	url := a.GetURL(urlID)

	if url == "" {
		http.Error(w, "Bad Url", 400)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
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
	code := a.AddURL(string(bodyBytes))
	shortenURL := a.GenerateShortenURL(code)

	w.WriteHeader(http.StatusCreated)

	_, writeError := w.Write([]byte(shortenURL))
	if writeError != nil {
		http.Error(w, "response body error", 400)
		return
	}

}
func (a *App) SaveURLJSONHandler(w http.ResponseWriter, r *http.Request) {
	var bodyBytes []byte
	var err error

	if r.Body != nil {
		bodyBytes, err = ioutil.ReadAll(r.Body)
		if err != nil || len(bodyBytes) == 0 {
			a.ShowJSONError(w, 400, "Only Json format requred in request body")
			return
		}
		defer r.Body.Close()
	}

	value := dataJSON{}
	if err := json.Unmarshal(bodyBytes, &value); err != nil {
		a.ShowJSONError(w, 400, "Only Json format requred in request body")
		return
	}

	code := a.AddURL(value.URL)
	shortenURL := a.GenerateShortenURL(code)
	shortenData := dataJSON{URL: shortenURL}
	shortenJSON, err := json.Marshal(shortenData)
	if err != nil {
		panic(err)
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
	return fmt.Sprintf("http://localhost:8080/%s", shortenCode)
}

func (a *App) ShowJSONError(w http.ResponseWriter, code int, message string) {
	dataJSON, err := json.Marshal(ErrorResponse{Code: code, Message: message})
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, writeError := w.Write(dataJSON)
	if writeError != nil {
		panic(writeError)
	}

}
