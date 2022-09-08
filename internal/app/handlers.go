package app

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rainset/shortener/internal/app/cookie"
	"io"
	"net/http"
)

func (a *App) NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/{id:[0-9a-zA-Z+]+}", a.GetURLHandler).Methods("GET")
	r.HandleFunc("/api/shorten", a.SaveURLJSONHandler).Methods("POST")
	r.HandleFunc("/api/user/urls", a.UserURLListHandler).Methods("GET")
	r.HandleFunc("/", a.SaveURLHandler).Methods("POST")
	return r
}

func (a *App) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urlValue := a.GetURL(vars["id"])
	if urlValue == "" {
		http.Error(w, "Bad Url", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, urlValue, http.StatusTemporaryRedirect)
}

func (a *App) UserURLListHandler(w http.ResponseWriter, r *http.Request) {

	type ListURL struct {
		ShortUrl    string `json:"short_url"`
		OriginalUrl string `json:"original_url"`
	}

	userId, err := cookie.Get(w, r, "userId")

	fmt.Println(userId, err)
	fmt.Println(a.userHistoryURLs)
	if err != nil || len(a.userHistoryURLs[userId]) == 0 {
		http.Error(w, "StatusNoContent", http.StatusNoContent)
		return
	}
	list := make([]ListURL, 0)

	for _, shortHashURL := range a.userHistoryURLs[userId] {

		originalUrl := a.urls[shortHashURL]

		if len(shortHashURL) == 0 || len(originalUrl) == 0 {
			continue
		}

		shortUrl := fmt.Sprintf("%s/%s", a.Config.ServerBaseURL, shortHashURL)
		list = append(list, ListURL{ShortUrl: shortUrl, OriginalUrl: originalUrl})
	}

	data, err := json.Marshal(list)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, writeError := w.Write([]byte(data))
	if writeError != nil {
		http.Error(w, "response body error", http.StatusBadRequest)
		return
	}
	//fmt.Println(list)
	//
	//data, err := json.Marshal(ErrorResponse{Code: code, Message: message})
	//if err != nil {
	//	panic(err)
	//}
	//w.Header().Set("Content-Type", "application/json")

}

func (a *App) SaveURLHandler(w http.ResponseWriter, r *http.Request) {
	var bodyBytes []byte
	var err error
	bodyBytes, err = readBodyBytes(r)

	if err != nil || len(bodyBytes) == 0 {
		http.Error(w, "Body reading error", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	code, err := a.AddURL(string(bodyBytes))
	if err != nil {
		http.Error(w, fmt.Sprintf("incorrect url format, code: %s body: %s", code, string(bodyBytes)), http.StatusBadRequest)
		return
	}

	generatedUserId := cookie.GenerateUniqueUserId()

	var cookieUserId string
	cookieUserId, err = cookie.Get(w, r, "userId")
	if err != nil {
		fmt.Println(err)
	}
	if len(cookieUserId) == 0 {
		cookie.Set(w, r, "userId", generatedUserId)
	}

	if err := a.AddUserHistoryURL(cookieUserId, code); err != nil {
		http.Error(w, "AddUserHistoryURL error", http.StatusBadRequest)
		return
	}

	shortenURL := a.GenerateShortenURL(code)

	//w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	_, writeError := w.Write([]byte(shortenURL))
	if writeError != nil {
		http.Error(w, "response body error", http.StatusBadRequest)
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
		bodyBytes, err = readBodyBytes(r)
		if err != nil || len(bodyBytes) == 0 {
			a.ShowJSONError(w, http.StatusBadRequest, "Only Json format required in request body")
			return
		}

		defer r.Body.Close()
	}

	value := ShortenRequest{}
	if err := json.Unmarshal(bodyBytes, &value); err != nil {
		a.ShowJSONError(w, http.StatusBadRequest, "Only Json format required in request body")
		return
	}

	code, err := a.AddURL(value.URL)

	fmt.Println("code:", code)
	fmt.Println("err:", err)

	if err != nil {
		http.Error(w, "incorrect url format", http.StatusBadRequest)
		return
	}

	shortenURL := a.GenerateShortenURL(code)
	shortenData := ShortenResponse{Result: shortenURL}
	shortenJSON, err := json.Marshal(shortenData)
	if err != nil {
		http.Error(w, "json response error", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, writeError := w.Write(shortenJSON)
	if writeError != nil {
		http.Error(w, "response body error", http.StatusBadRequest)
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

func readBodyBytes(r *http.Request) ([]byte, error) {
	// Read body
	bodyBytes, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		return nil, readErr
	}
	defer r.Body.Close()

	// GZIP decode
	if len(r.Header["Content-Encoding"]) > 0 && r.Header["Content-Encoding"][0] == "gzip" {
		r, gzErr := gzip.NewReader(io.NopCloser(bytes.NewBuffer(bodyBytes)))
		if gzErr != nil {
			return nil, gzErr
		}
		defer r.Close()

		bb, err2 := io.ReadAll(r)
		if err2 != nil {
			return nil, err2
		}
		return bb, nil
	} else {
		// Not compressed
		return bodyBytes, nil
	}
}
