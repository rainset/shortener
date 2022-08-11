package app

import (
	"crypto/md5"
	"encoding/hex"
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

func New() *App {
	return &App{urls: make(map[string]string)}
}

func (a *App) AddURL(value string) string {
	a.Lock()
	defer a.Unlock()
	binHash := md5.Sum([]byte(value))
	hash := hex.EncodeToString(binHash[:])
	a.urls[hash] = value
	return hash
}

func (a *App) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urlID := vars["id"]
	url := a.urls[urlID]

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
	shortenURL := fmt.Sprintf("http://localhost:8080/%s", code)

	w.WriteHeader(http.StatusCreated)

	_, writeError := w.Write([]byte(shortenURL))
	if writeError != nil {
		http.Error(w, "response body error", 400)
		return
	}

}
