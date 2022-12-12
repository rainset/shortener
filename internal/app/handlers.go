// Package app - http handlers

package app

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/rainset/shortener/internal/helper"
	queue "github.com/rainset/shortener/internal/queue"
	"github.com/rainset/shortener/internal/storage"
)

// NewRouter регистрирует все хендлеры,middleware, сессии и cookie
func (a *App) NewRouter() *gin.Engine {

	r := gin.Default()
	//r.Use(gzip_gin.Gzip(gzip_gin.DefaultCompression))

	store := cookie.NewStore([]byte(a.Config.CookieHashKey), []byte(a.Config.CookieBlockKey))
	store.Options(sessions.Options{MaxAge: 3600})
	r.Use(sessions.Sessions("sessid", store))

	r.GET("/ping", a.PingHandler)
	r.GET("/:id", a.GetURLHandler)
	r.POST("/api/shorten/batch", a.SaveURLBatchJSONHandler)
	r.POST("/api/shorten", a.SaveURLJSONHandler)
	r.DELETE("/api/user/urls", a.DeleteUserBatchURLHandler)
	r.GET("/api/user/urls", a.UserURLListHandler)
	r.POST("/", a.SaveURLHandler)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound})
	})
	return r
}

// PingHandler хендлер GET /ping
// проверка доступности сервиса
func (a *App) PingHandler(c *gin.Context) {
	err := a.s.Ping()
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "err": "ping error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK})
}

// GetURLHandler хендлер POST /:id
// перед по короткой ссылке
func (a *App) GetURLHandler(c *gin.Context) {

	id := c.Param("id")
	resultURL, err := a.s.GetURL(id)

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if resultURL.Deleted == 1 {
		c.AbortWithStatus(http.StatusGone)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, resultURL.Original)

}

// UserURLListHandler хендлер GET /api/user/urls
// просмотр списка добавленных ссылок пользователем
func (a *App) UserURLListHandler(c *gin.Context) {

	type ListURL struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	var userHistoryURLs []storage.ResultHistoryURL
	var err error
	var userID string

	ss := sessions.Default(c)
	session, ok := ss.Get("sessid").(Session)
	if ok {
		userID = session.UserID
		userHistoryURLs, err = a.s.GetListUserHistoryURL(userID)
	}

	if err != nil || len(userID) == 0 || len(userHistoryURLs) == 0 {
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusNoContent)
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
	c.JSON(http.StatusOK, list)
}

// SaveURLHandler хендлер POST /
func (a *App) SaveURLHandler(c *gin.Context) {

	var cookieuserID string
	var bodyBytes []byte
	var err error

	ss := sessions.Default(c)
	session, ok := ss.Get("sessid").(Session)
	if ok {
		cookieuserID = session.UserID
	} else {
		genID := helper.GenerateUniqueuserID()
		ss.Set("sessid", Session{UserID: genID})
		_ = ss.Save()
		cookieuserID = genID
	}

	bodyBytes, err = readBodyBytes(c)

	if err != nil || len(bodyBytes) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "body error"})
		return
	}

	//urlValue, err := url.ParseRequestURI(string(bodyBytes))
	urlValue := string(bodyBytes)
	if err != nil {
		return
	}
	hash := helper.GenerateToken(8)

	var isDBExist bool
	err = a.s.AddURL(hash, urlValue)
	if err != nil {
		fmt.Println("a.s.AddURL:", err)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				isDBExist = true
				hash, err = a.s.GetByOriginalURL(urlValue)
				if err != nil {
					return
				}
			}
		}
	}

	if err != nil && !isDBExist {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("incorrect url format, hash: %s body: %s", hash, urlValue)})
		return
	}

	shortenURL := a.GenerateShortenURL(hash)

	err = a.s.AddUserHistoryURL(cookieuserID, hash)

	if isDBExist {
		c.String(http.StatusConflict, "%s", shortenURL)
		return
	} else if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	} else {
		c.String(http.StatusCreated, "%s", shortenURL)
		return
	}
}

// SaveURLJSONHandler хендлер POST /api/shorten
func (a *App) SaveURLJSONHandler(c *gin.Context) {

	var err error

	type requestBody struct {
		URL string `json:"url"`
	}
	type responseBody struct {
		Result string `json:"result"`
	}
	value := requestBody{}

	err = c.BindJSON(&value)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Only Json format required in request body"})
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
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	shortenURL := a.GenerateShortenURL(hash)
	shortenData := responseBody{Result: shortenURL}

	if isDBExist {
		c.JSON(http.StatusConflict, shortenData)
		return
	}

	c.JSON(http.StatusCreated, shortenData)
}

// SaveURLBatchJSONHandler хендлер POST /api/shorten/batch
func (a *App) SaveURLBatchJSONHandler(c *gin.Context) {

	type ShortenBatchRequest struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	type ShortenBatchResponse struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}

	var err error

	requestBody := make([]ShortenBatchRequest, 0)

	err = c.BindJSON(&requestBody)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Only Json format required in request body"})
		return
	}

	batchURLs := make([]storage.BatchUrls, 0)
	for _, v := range requestBody {
		batchURLs = append(batchURLs, storage.BatchUrls{CorrelationID: v.CorrelationID, OriginalURL: v.OriginalURL})
	}

	result, err := a.s.AddBatchURL(batchURLs)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "db save error", "err": err})
		return
	}

	var response []ShortenBatchResponse
	for _, v := range result {
		response = append(response, ShortenBatchResponse{ShortURL: a.GenerateShortenURL(v.Hash), CorrelationID: v.CorrelationID})
	}

	c.JSON(http.StatusCreated, response)
}

// DeleteUserBatchURLHandler хендлер DELETE /api/user/urls
func (a *App) DeleteUserBatchURLHandler(c *gin.Context) {
	var err error
	var hashes []string

	ss := sessions.Default(c)
	session, ok := ss.Get("sessid").(Session)
	if ok {
		err = c.BindJSON(&hashes)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Only Json format required in request body"})
			return
		}
		a.Queue.Push(&queue.Task{CookieID: session.UserID, Hashes: hashes})
	}

	c.Status(http.StatusAccepted)
}

// readBodyBytes распаковывает тело запроса в формате gzip
func readBodyBytes(c *gin.Context) ([]byte, error) {
	bodyBytes, readErr := io.ReadAll(c.Request.Body)
	if readErr != nil {
		return nil, readErr
	}

	if strings.Contains(c.GetHeader("Content-Encoding"), "gzip") {
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
	}
	// Not compressed
	return bodyBytes, nil
}
