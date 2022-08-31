package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
)

func (a *App) NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/{id:[0-9a-z]+}", a.GetURLHandler).Methods("GET")
	r.HandleFunc("/api/shorten", a.SaveURLJSONHandler).Methods("POST")
	r.HandleFunc("/", a.SaveURLHandler).Methods("POST")

	//handlers.CompressHandler(r)

	return r
}

func (a *App) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urlValue := a.GetURL(vars["id"])
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

		var reader io.Reader

		reader = r.Body

		//switch r.Header.Get("Content-Encoding") {
		//case "gzip":
		//	reader, err := gzip.NewReader(r.Body)
		//	if err != nil {
		//		http.Error(w, err.Error(), http.StatusInternalServerError)
		//		return
		//	}
		//	defer reader.Close()
		//default:
		//	fmt.Println("default")
		//	reader = r.Body
		//}

		bodyBytes, err = ioutil.ReadAll(reader)
		if err != nil || len(bodyBytes) == 0 {

			fmt.Println(err, bodyBytes, string(bodyBytes))

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

	//w.Header().Set("Content-Type", "text/plain")
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
