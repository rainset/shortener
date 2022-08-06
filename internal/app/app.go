package app

import (
	"net/http"
	"strings"
)

type App struct {
	urls map[string]string
}

func New() *App {
	return &App{urls: make(map[string]string)}
}

func (a App) AddUrl(code, value string) {
	a.urls[code] = value
}

func (a App) GetUrl(code string) string {
	return a.urls[code]
}

func (a *App) RouteHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		p := strings.Split(r.URL.Path, "/")[1:]
		urlId := p[0]
		if urlId == "" {
			http.Error(w, "Bad Url", 400)
			return
		}
		url := "http://cehme.ru/ces4suo/a1dez"
		http.Redirect(w, r, url, http.StatusMovedPermanently)

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			return
		}
		//url := r.FormValue("url")
		w.WriteHeader(http.StatusCreated)
		shortedUrl := a.urls["link1"]
		w.Write([]byte(shortedUrl))

	default:
		http.Error(w, "Sorry, only GET and POST methods are supported.", 400)
		return
	}
}
