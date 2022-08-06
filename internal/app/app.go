package app

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type App struct {
	urls map[string]string
}

func New() *App {
	return &App{urls: make(map[string]string)}
}

func (a App) AddUrl(value string) string {
	seconds := time.Now().Second()
	randInt := rand.Intn(9999999)
	code := fmt.Sprintf("%d%d", randInt, seconds)
	a.urls[code] = value
	return code
}

func (a App) GetUrl(code string) string {
	return a.urls[code]
}

func (a App) RouteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p := strings.Split(r.URL.Path, "/")[1:]
		urlId := p[0]
		if urlId == "" {
			http.Error(w, "Bad Url", 400)
			return
		}

		url, ok := a.urls[urlId]

		if !ok {
			http.Error(w, "Bad Url", 400)
			return
		}

		http.Redirect(w, r, url, http.StatusMovedPermanently)

	case http.MethodPost:
		var bodyBytes []byte
		var err error

		if r.Body != nil {
			bodyBytes, err = ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Printf("Body reading error: %v", err)
				return
			}
			defer r.Body.Close()
		}
		code := a.AddUrl(string(bodyBytes))
		shortenUrl := fmt.Sprintf("http://localhost:8080/%s", code)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortenUrl))

	default:
		http.Error(w, "Sorry, only GET and POST methods are supported.", 400)
		return
	}
}
