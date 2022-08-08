package app

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type App struct {
	urls map[string]string
}

func New() *App {
	return &App{urls: make(map[string]string)}
}

func (a *App) AddUrl(value string) string {
	binHash := md5.Sum([]byte(value))
	hash := hex.EncodeToString(binHash[:])
	a.urls[hash] = value
	return hash
}

func (a *App) GetUrlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urlId := vars["id"]
	url := a.urls[urlId]

	if url == "" {
		http.Error(w, "Bad Url", 400)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (a *App) SaveUrlHandler(w http.ResponseWriter, r *http.Request) {
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
	code := a.AddUrl(string(bodyBytes))

	shortenUrl := fmt.Sprintf("http://localhost:8080/%s", code)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortenUrl))
}
