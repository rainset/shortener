package app

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/rainset/shortener/internal/helper"
	"github.com/rainset/shortener/internal/storage"
	"io"
	"net/http"
	"net/url"
	"sync"
)

func (a *App) NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/ping", a.PingHandler).Methods("GET")
	r.HandleFunc("/{id:[0-9a-zA-Z+]+}", a.GetURLHandler).Methods("GET")
	r.HandleFunc("/api/shorten/batch", a.SaveURLBatchJSONHandler).Methods("POST")
	r.HandleFunc("/api/shorten", a.SaveURLJSONHandler).Methods("POST")
	r.HandleFunc("/api/user/urls", a.DeleteBatchURLHandler).Methods("DELETE")
	r.HandleFunc("/api/user/urls", a.UserURLListHandler).Methods("GET")
	r.HandleFunc("/", a.SaveURLHandler).Methods("POST")
	return r
}

func (a *App) PingHandler(w http.ResponseWriter, _ *http.Request) {
	err := a.s.Ping()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *App) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resultURL, err := a.s.GetURL(vars["id"])

	if err != nil {
		http.Error(w, "Bad Url", http.StatusBadRequest)
		return
	}

	if resultURL.Deleted == 1 {
		http.Error(w, "Bad Url", http.StatusGone)
		return
	}

	http.Redirect(w, r, resultURL.Original, http.StatusTemporaryRedirect)
}

func (a *App) UserURLListHandler(w http.ResponseWriter, r *http.Request) {

	type ListURL struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	userID, _ := a.cookie.Get(w, r, "userID")
	userHistoryURLs, err := a.s.GetListUserHistoryURL(userID)

	fmt.Println(userID, err)
	fmt.Println("UserURLListHandler", userHistoryURLs)
	if err != nil || len(userID) == 0 || len(userHistoryURLs) == 0 {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)

		return
	}
	list := make([]ListURL, 0)

	for _, item := range userHistoryURLs {
		if len(item.Hash) == 0 || len(item.CookieID) == 0 {
			continue
		}
		ShortURL := fmt.Sprintf("%s/%s", a.Config.ServerBaseURL, item.Hash)
		list = append(list, ListURL{ShortURL: ShortURL, OriginalURL: item.Original})
	}

	data, err := json.Marshal(list)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, writeError := w.Write(data)
	if writeError != nil {
		http.Error(w, "response body error", http.StatusBadRequest)
		return
	}
}

func (a *App) SaveURLHandler(w http.ResponseWriter, r *http.Request) {

	var bodyBytes []byte
	var err error
	bodyBytes, err = readBodyBytes(r)

	if err != nil || len(bodyBytes) == 0 {
		http.Error(w, "Body reading error", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	urlValue, err := url.ParseRequestURI(string(bodyBytes))
	if err != nil {
		return
	}
	hash := helper.GenerateToken(8)

	var isDBExist bool
	err = a.s.AddURL(hash, urlValue.String())
	if err != nil {
		fmt.Println("a.s.AddURL:", err)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				isDBExist = true
				hash, err = a.s.GetByOriginalURL(urlValue.String())
				if err != nil {
					return
				}
			}
		}
	}

	if err != nil && !isDBExist {
		http.Error(w, fmt.Sprintf("incorrect url format, hash: %s body: %s", hash, urlValue.String()), http.StatusBadRequest)
		return
	}

	shortenURL := a.GenerateShortenURL(hash)
	generateduserID := helper.GenerateUniqueuserID()
	cookieuserID, err := a.cookie.Get(w, r, "userID")
	if err != nil {
		fmt.Println(err)
	}
	if len(cookieuserID) == 0 {
		a.cookie.Set(w, r, "userID", generateduserID)
		cookieuserID = generateduserID
	}

	err = a.s.AddUserHistoryURL(cookieuserID, hash)

	if isDBExist {
		w.WriteHeader(http.StatusConflict)
	} else if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
	}

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

		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(r.Body)
	}
	fmt.Println("")
	value := ShortenRequest{}
	if err := json.Unmarshal(bodyBytes, &value); err != nil {
		a.ShowJSONError(w, http.StatusBadRequest, "Only Json format required in request body")
		return
	}

	hash := helper.GenerateToken(8)

	var isDBExist bool
	err = a.s.AddURL(hash, value.URL)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				isDBExist = true
				hash, err = a.s.GetByOriginalURL(value.URL)
				if err != nil {
					return
				}
			}
		}
	}

	if err != nil {
		http.Error(w, "incorrect url format", http.StatusBadRequest)
		return
	}

	shortenURL := a.GenerateShortenURL(hash)
	shortenData := ShortenResponse{Result: shortenURL}
	shortenJSON, err := json.Marshal(shortenData)

	w.Header().Set("Content-Type", "application/json")

	if isDBExist {
		w.WriteHeader(http.StatusConflict)
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	_, writeError := w.Write(shortenJSON)
	if writeError != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func (a *App) SaveURLBatchJSONHandler(w http.ResponseWriter, r *http.Request) {

	type ShortenBatchRequest struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	type ShortenBatchResponse struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}

	var bodyBytes []byte
	var err error

	if r.Body != nil {
		bodyBytes, err = readBodyBytes(r)
		if err != nil || len(bodyBytes) == 0 {
			a.ShowJSONError(w, http.StatusBadRequest, "Only Json format required in request body")
			return
		}

		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(r.Body)
	}

	shortenBatchRequestList := make([]ShortenBatchRequest, 0)
	batchURLs := make([]storage.BatchUrls, 0)
	//var value []interface{}
	if err = json.Unmarshal(bodyBytes, &shortenBatchRequestList); err != nil {

		//fmt.Println(err)
		//fmt.Println(shortenBatchRequestList)
		a.ShowJSONError(w, http.StatusBadRequest, "json decode error")
		return
	}
	fmt.Println(shortenBatchRequestList)

	for _, v := range shortenBatchRequestList {
		batchURLs = append(batchURLs, storage.BatchUrls{CorrelationID: v.CorrelationID, OriginalURL: v.OriginalURL})
	}

	result, err := a.s.AddBatchURL(batchURLs)
	if err != nil {
		fmt.Println(err)
		a.ShowJSONError(w, http.StatusBadRequest, "db save error")
		return
	}

	var response []ShortenBatchResponse
	for _, v := range result {
		response = append(response, ShortenBatchResponse{ShortURL: a.GenerateShortenURL(v.Hash), CorrelationID: v.CorrelationID})
	}

	shortenJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "json response error", http.StatusBadRequest)
		return
	}
	fmt.Println(shortenJSON, string(shortenJSON))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, writeError := w.Write(shortenJSON)
	if writeError != nil {
		http.Error(w, "response body error", http.StatusBadRequest)
		return
	}
}

func (a *App) DeleteBatchURLHandler(w http.ResponseWriter, r *http.Request) {

	var bodyBytes []byte
	var err error

	if r.Body != nil {
		bodyBytes, err = readBodyBytes(r)
		if err != nil || len(bodyBytes) == 0 {
			a.ShowJSONError(w, http.StatusBadRequest, "Only Json format required in request body")
			return
		}

		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(r.Body)
	}

	var mapHashes []string

	if err := json.Unmarshal(bodyBytes, &mapHashes); err != nil {
		a.ShowJSONError(w, http.StatusBadRequest, "json decode error")
		return
	}

	chunkSize := 50
	var chunks [][]string
	for i := 0; i < len(mapHashes); i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(mapHashes) {
			end = len(mapHashes)
		}

		chunks = append(chunks, mapHashes[i:end])
	}

	cookieID, _ := a.cookie.Get(w, r, "userID")
	var wg sync.WaitGroup
	for _, v := range chunks {
		wg.Add(1)
		go func(a *App, cookieID string, v []string) {
			fmt.Println(v, cookieID)
			_ = a.s.DeleteBatchURL(cookieID, v)
			wg.Done()
		}(a, cookieID, v)
	}

	wg.Wait() // ждём все горутины

	w.WriteHeader(http.StatusAccepted)

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
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	// GZIP decode
	if len(r.Header["Content-Encoding"]) > 0 && r.Header["Content-Encoding"][0] == "gzip" {
		r, gzErr := gzip.NewReader(io.NopCloser(bytes.NewBuffer(bodyBytes)))
		if gzErr != nil {
			return nil, gzErr
		}
		defer func(r *gzip.Reader) {
			_ = r.Close()
		}(r)

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
