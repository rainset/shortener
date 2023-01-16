// Package app - http handlers

package app

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"net/netip"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/rainset/shortener/internal/helper"
	queue "github.com/rainset/shortener/internal/queue"
)

type ShortenerHTTPServer struct {
	a *App
}

// NewRouter регистрирует все хендлеры,middleware, сессии и cookie
func (s *ShortenerHTTPServer) NewRouter() *gin.Engine {

	r := gin.Default()

	store := cookie.NewStore([]byte(s.a.Config.CookieHashKey), []byte(s.a.Config.CookieBlockKey))
	store.Options(sessions.Options{MaxAge: 3600})
	r.Use(sessions.Sessions("sessid", store))

	r.GET("/ping", s.PingHandler)
	r.GET("/:id", s.GetURLHandler)
	r.POST("/api/shorten/batch", s.SaveURLBatchJSONHandler)
	r.POST("/api/shorten", s.SaveURLJSONHandler)
	r.DELETE("/api/user/urls", s.DeleteUserBatchURLHandler)
	r.GET("/api/user/urls", s.UserURLListHandler)
	r.GET("/api/internal/stats", s.SubnetMiddleware, s.StatsHandler)
	r.POST("/", s.SaveURLHandler)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound})
	})
	return r
}

// SubnetMiddleware
// Middleware проверяет принадлежность IP адреса клиента на совпадение разрешенной маски безклассовой адресации (Classless Inter-Domain Routing, CIDR)
func (s *ShortenerHTTPServer) SubnetMiddleware(c *gin.Context) {

	network, err := netip.ParsePrefix(s.a.Config.TrustedSubnet)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden})
	}

	ip, err := netip.ParseAddr(c.ClientIP())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden})
	}

	b := network.Contains(ip)
	log.Println("network.Contains:", c.ClientIP(), b) // true

	if !b {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden})
	}

	c.Next()
}

// StatsHandler хендлер GET /api/internal/stats
// статистика сервиса, доступ по маске IP
func (s *ShortenerHTTPServer) StatsHandler(c *gin.Context) {
	stats, err := s.a.GetStats()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": stats.Users, "urls": stats.Urls})
}

// PingHandler хендлер GET /ping
// проверка доступности сервиса
func (s *ShortenerHTTPServer) PingHandler(c *gin.Context) {
	err := s.a.Ping()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK})
}

// GetURLHandler хендлер POST /:id
// перед по короткой ссылке
func (s *ShortenerHTTPServer) GetURLHandler(c *gin.Context) {
	id := c.Param("id")
	resultURL, err := s.a.GetURL(id)
	if err != nil {
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
func (s *ShortenerHTTPServer) UserURLListHandler(c *gin.Context) {

	ss := sessions.Default(c)
	session, ok := ss.Get("sessid").(Session)
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userID := session.UserID
	list, err := s.a.GetListUserHistoryURL(userID)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	if len(userID) == 0 || len(list) == 0 {
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, list)
}

// SaveURLHandler хендлер POST /
func (s *ShortenerHTTPServer) SaveURLHandler(c *gin.Context) {

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

	urlValue := string(bodyBytes)
	addURLResult, err := s.a.AddURL(urlValue)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = s.a.AddUserHistoryURL(cookieuserID, addURLResult.Hash)

	if addURLResult.Exists {
		c.String(http.StatusConflict, "%s", addURLResult.ShortURL)
		return
	} else if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	} else {
		c.String(http.StatusCreated, "%s", addURLResult.ShortURL)
		return
	}
}

// SaveURLJSONHandler хендлер POST /api/shorten
func (s *ShortenerHTTPServer) SaveURLJSONHandler(c *gin.Context) {

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
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	addURLResult, err := s.a.AddURL(value.URL)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	shortenData := responseBody{Result: addURLResult.ShortURL}

	if addURLResult.Exists {
		c.JSON(http.StatusConflict, shortenData)
		return
	}

	c.JSON(http.StatusCreated, shortenData)
}

// SaveURLBatchJSONHandler хендлер POST /api/shorten/batch
func (s *ShortenerHTTPServer) SaveURLBatchJSONHandler(c *gin.Context) {

	requestBody := make([]AddURLBatchRequest, 0)
	err := c.BindJSON(&requestBody)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	response, err := s.a.AddBatchURL(requestBody)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// DeleteUserBatchURLHandler хендлер DELETE /api/user/urls
func (s *ShortenerHTTPServer) DeleteUserBatchURLHandler(c *gin.Context) {
	var err error
	var hashes []string

	ss := sessions.Default(c)
	session, ok := ss.Get("sessid").(Session)
	if ok {
		err = c.BindJSON(&hashes)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		s.a.Queue.Push(&queue.Task{CookieID: session.UserID, Hashes: hashes})
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
